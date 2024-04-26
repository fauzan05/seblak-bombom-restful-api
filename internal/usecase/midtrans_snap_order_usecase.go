package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MidtransSnapOrderUseCase struct {
	Log                         *logrus.Logger
	DB                          *gorm.DB
	Validate                    *validator.Validate
	SnapClient                  *snap.Client
	OrderRepository             *repository.OrderRepository
	MidtransSnapOrderRepository *repository.MidtransSnapOrderRepository
}

func NewMidtransSnapOrderUseCase(log *logrus.Logger, validate *validator.Validate, orderRepository *repository.OrderRepository,
	snapClient *snap.Client, db *gorm.DB, midtransSnapOrderRepository *repository.MidtransSnapOrderRepository) *MidtransSnapOrderUseCase {
	return &MidtransSnapOrderUseCase{
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		SnapClient:                  snapClient,
		DB:                          db,
		MidtransSnapOrderRepository: midtransSnapOrderRepository,
	}
}

func (c *MidtransSnapOrderUseCase) Add(ctx context.Context, request *model.CreateMidtransSnapOrderRequest) (*model.MidtransSnapOrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newSnapClient := snap.Client{}
	newSnapClient.New(c.SnapClient.ServerKey, c.SnapClient.Env)

	// temukan data order berdasarkan invoice dari request
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if selectedOrder.PaymentMethod != helper.ONLINE {
		c.Log.Warnf("This order has an onsite payment method : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var midtransItemDetails []midtrans.ItemDetails
	for _, product := range selectedOrder.OrderProducts {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID:    strconv.Itoa(int(product.ID)),
			Qty:   int32(product.Quantity),
			Price: int64(product.Price),
			Name:  product.ProductName,
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}
	// cek apakah ada biaya pengiriman, jika ada maka tambahkan ke midtransItemDetails agar value-nya sama
	if selectedOrder.DeliveryCost != 0 {
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

	orderId := int(selectedOrder.ID)
	midtransRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(orderId),
			GrossAmt: int64(selectedOrder.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: selectedOrder.FirstName,
			LName: selectedOrder.LastName,
			Email: selectedOrder.Email,
			Phone: selectedOrder.Phone,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &midtransItemDetails,
	}

	snapResponse, snapErr := newSnapClient.CreateTransaction(midtransRequest)
	if snapErr != nil {
		c.Log.Warnf("Failed to create new midtrans transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newMidtransSnapOrder := new(entity.MidtransSnapOrder)
	newMidtransSnapOrder.OrderId = request.OrderId
	newMidtransSnapOrder.Token = snapResponse.Token
	newMidtransSnapOrder.RedirectUrl = snapResponse.RedirectURL
	if err := c.MidtransSnapOrderRepository.Create(tx, newMidtransSnapOrder); err != nil {
		c.Log.Warnf("Failed to create new midtrans snap order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.MidtransSnapOrderToResponse(newMidtransSnapOrder), nil
}

func (c *MidtransSnapOrderUseCase) Get(ctx context.Context, request *model.GetMidtransSnapOrderRequest) (*model.MidtransSnapOrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newMidtransSnapOrder := new(entity.MidtransSnapOrder)
	if err := c.MidtransSnapOrderRepository.FindMidtransSnapOrderByOrderId(tx, newMidtransSnapOrder, request.OrderId); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.MidtransSnapOrderToResponse(newMidtransSnapOrder), nil
}