package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MidtransCoreAPIOrderUseCase struct {
	Log                            *logrus.Logger
	DB                             *gorm.DB
	Validate                       *validator.Validate
	CoreAPIClient                  *coreapi.Client
	OrderRepository                *repository.OrderRepository
	ProductRepository              *repository.ProductRepository
	MidtransCoreAPIOrderRepository *repository.MidtransCoreAPIOrderRepository
}

var layoutWithoutZone = "2006-01-02 15:04:05"

func NewMidtransCoreAPIOrderUseCase(log *logrus.Logger, validate *validator.Validate, orderRepository *repository.OrderRepository,
	coreAPIClient *coreapi.Client, db *gorm.DB, midtransCoreAPiOrderRepository *repository.MidtransCoreAPIOrderRepository,
	productRepository *repository.ProductRepository) *MidtransCoreAPIOrderUseCase {
	return &MidtransCoreAPIOrderUseCase{
		Log:                            log,
		Validate:                       validate,
		OrderRepository:                orderRepository,
		ProductRepository:              productRepository,
		CoreAPIClient:                  coreAPIClient,
		DB:                             db,
		MidtransCoreAPIOrderRepository: midtransCoreAPiOrderRepository,
	}
}

func (c *MidtransCoreAPIOrderUseCase) Add(ctx context.Context, request *model.CreateMidtransCoreAPIOrderRequest) (*model.MidtransCoreAPIOrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCoreAPIClient := coreapi.Client{}
	newCoreAPIClient.New(c.CoreAPIClient.ServerKey, c.CoreAPIClient.Env)

	// temukan data order berdasarkan invoice dari request
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	var midtransItemDetails []midtrans.ItemDetails
	for _, product := range selectedOrder.OrderProducts {
		newMidtransItemDetail := midtrans.ItemDetails{
			ID:    strconv.Itoa(int(product.ID)),
			Qty:   int32(product.Quantity),
			Price: int64(product.Price),
			Name:  product.ProductName,
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemDetail)
	}
	// cek apakah ada biaya pengiriman, jika ada maka tambahkan ke midtransItemDetails agar value-nya sama
	if selectedOrder.DeliveryCost > 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID:    "1",
			Qty:   1,
			Price: int64(selectedOrder.DeliveryCost),
			Name:  "Delivery Cost",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}
	// cek juga apakah terdapat diskon, jika ada maka tambahkan ke ItemDetails
	if selectedOrder.TotalDiscount > 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID:    "2",
			Qty:   1,
			Price: -int64(selectedOrder.TotalDiscount),
			Name:  "Discount",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}

	var paymentType coreapi.CoreapiPaymentType
	if selectedOrder.PaymentMethod == "" {
		c.Log.Warnf("Payment method is not set : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if selectedOrder.PaymentMethod == helper.GOPAY {
		paymentType = coreapi.PaymentTypeGopay
	} else if selectedOrder.PaymentMethod == helper.QRIS {
		paymentType = coreapi.PaymentTypeQris
	}

	midtransRequest := coreapi.ChargeReq{
		PaymentType: paymentType,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(selectedOrder.ID)),
			GrossAmt: int64(selectedOrder.Amount),
		},
		Items: &midtransItemDetails,
	}

	coreApiResponse, coreApiErr := newCoreAPIClient.ChargeTransaction(&midtransRequest)
	if coreApiErr != nil {
		c.Log.Warnf("Failed to create new midtrans transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create a new transaction! Please try again later!")
	}

	newMidtransCoreAPIOrder := new(entity.MidtransCoreAPIOrder)
	newMidtransCoreAPIOrder.OrderId = selectedOrder.ID
	newMidtransCoreAPIOrder.MidtransOrderId = coreApiResponse.OrderID
	newMidtransCoreAPIOrder.StatusCode = coreApiResponse.StatusCode
	newMidtransCoreAPIOrder.StatusMessage = coreApiResponse.StatusMessage
	newMidtransCoreAPIOrder.TransactionId = coreApiResponse.TransactionID
	grossAmount64, err := strconv.ParseFloat(coreApiResponse.GrossAmount, 64)
	if err != nil {
		c.Log.Warnf("Failed to parse gross amount into float32 : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newMidtransCoreAPIOrder.GrossAmount = float32(grossAmount64)
	newMidtransCoreAPIOrder.Currency = coreApiResponse.Currency
	newMidtransCoreAPIOrder.PaymentType = coreApiResponse.PaymentType
	parseExpiryTime, err := time.ParseInLocation(layoutWithoutZone, coreApiResponse.ExpiryTime, time.Local)
	if err != nil {
		c.Log.Warnf("Failed to parse expiry time into standart format : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newMidtransCoreAPIOrder.ExpiryTime = parseExpiryTime

	parseTransactionTime, err := time.ParseInLocation(layoutWithoutZone, coreApiResponse.TransactionTime, time.Local)
	if err != nil {
		c.Log.Warnf("Failed to parse transaction time into standart format : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newMidtransCoreAPIOrder.TransactionTime = parseTransactionTime
	newMidtransCoreAPIOrder.TransactionStatus = helper.TransactionStatus(coreApiResponse.TransactionStatus)
	newMidtransCoreAPIOrder.FraudStatus = coreApiResponse.FraudStatus
	for _, action := range coreApiResponse.Actions {
		newMidtransCoreAPIOrder.Actions = append(newMidtransCoreAPIOrder.Actions, entity.Action{
			Name:   action.Name,
			Method: helper.RequestMethod(action.Method),
			URL:    action.URL,
		})
	}

	if err := c.MidtransCoreAPIOrderRepository.Create(tx, newMidtransCoreAPIOrder); err != nil {
		c.Log.Warnf("Failed to create new midtrans core api order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to store data into database!")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.MidtransCoreAPIToResponse(newMidtransCoreAPIOrder), nil
}

func (c *MidtransCoreAPIOrderUseCase) Get(ctx context.Context, request *model.GetMidtransCoreAPIOrderRequest) (*model.MidtransCoreAPIOrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	selectedMidtransOrder := new(entity.MidtransCoreAPIOrder)
	if err := c.MidtransCoreAPIOrderRepository.FindMidtransCoreAPIOrderByOrderId(tx, selectedMidtransOrder, request.OrderId); err != nil {
		c.Log.Warnf("Failed to find midtrans order by order id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newCoreAPIClient := coreapi.Client{}
	newCoreAPIClient.New(c.CoreAPIClient.ServerKey, c.CoreAPIClient.Env)
	parseToUint := strconv.FormatUint(selectedMidtransOrder.OrderId, 10)
	coreAPIResponse, err := newCoreAPIClient.CheckTransaction(parseToUint)

	if selectedMidtransOrder.TransactionStatus != helper.TransactionStatus(coreAPIResponse.TransactionStatus) {
		// update statusnya
		selectedMidtransOrder.TransactionStatus = helper.TransactionStatus(coreAPIResponse.TransactionStatus)
		if err := c.MidtransCoreAPIOrderRepository.Update(tx, selectedMidtransOrder); err != nil {
			c.Log.Warnf("Failed to update midtrans transaction status : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}
	
	return converter.MidtransCoreAPIToResponse(selectedMidtransOrder), nil
}
