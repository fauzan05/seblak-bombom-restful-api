package usecase

import (
	"encoding/json"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"time"

	"seblak-bombom-restful-api/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payment_request"
	"gorm.io/gorm"
)

type XenditTransactionQRCodeUseCase struct {
	DB                          *gorm.DB
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	XenditClient                *xendit.APIClient
	OrderRepository             *repository.OrderRepository
	XenditTransactionRepository *repository.XenditTransctionRepository
}

func NewXenditTransactionQRCodeUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditClient *xendit.APIClient) *XenditTransactionQRCodeUseCase {
	return &XenditTransactionQRCodeUseCase{
		DB:                          db,
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		XenditTransactionRepository: xenditTransactionRepository,
		XenditClient:                xenditClient,
	}
}

func (c *XenditTransactionQRCodeUseCase) Add(ctx *fiber.Ctx, request *model.CreateXenditTransaction) (*model.XenditTransactionResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// get order id
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	paymentRequestBasketItems := new([]payment_request.PaymentRequestBasketItem)
	for _, product := range selectedOrder.OrderProducts {
		refId := strconv.FormatUint(product.ProductId, 10)
		itemType := string(helper.ITEM_TYPE_PHYSICAL_PRODUCT)
		paymentRequestBasketItem := &payment_request.PaymentRequestBasketItem{
			ReferenceId: &refId,
			Name:        product.ProductName,
			Currency:    string(payment_request.PAYMENTREQUESTCURRENCY_IDR),
			Quantity:    float64(product.Quantity),
			Price:       float64(product.Price),
			Category:    product.Category,
			Type:        &itemType,
		}
		*paymentRequestBasketItems = append(*paymentRequestBasketItems, *paymentRequestBasketItem)
	}

	// cek apakah ada biaya pengiriman
	if selectedOrder.DeliveryCost > 0 {
		refId := fmt.Sprintf("DELIVERY/%s", strconv.FormatUint(selectedOrder.ID, 10))
		itemType := string(helper.ITEM_TYPE_DELIVERY_FEE)
		paymentRequestBasketItem := &payment_request.PaymentRequestBasketItem{
			ReferenceId: &refId,
			Name:        "Delivery Cost",
			Currency:    string(payment_request.PAYMENTREQUESTCURRENCY_IDR),
			Quantity:    1,
			Price:       float64(selectedOrder.DeliveryCost),
			Category:    "delivery",
			Type:        &itemType,
		}
		*paymentRequestBasketItems = append(*paymentRequestBasketItems, *paymentRequestBasketItem)
	}

	if selectedOrder.TotalDiscount > 0 {
		refId := fmt.Sprintf("DISCOUNT/%s", strconv.FormatUint(selectedOrder.ID, 10))
		itemType := string(helper.ITEM_TYPE_DISCOUNT)
		paymentRequestBasketItem := &payment_request.PaymentRequestBasketItem{
			ReferenceId: &refId,
			Name:        "Discount",
			Currency:    string(payment_request.PAYMENTREQUESTCURRENCY_IDR),
			Quantity:    1,
			Price:       float64(selectedOrder.TotalDiscount),
			Category:    "discount",
			Type:        &itemType,
		}
		*paymentRequestBasketItems = append(*paymentRequestBasketItems, *paymentRequestBasketItem)
	}

	amountFloat64 := float64(selectedOrder.Amount)
	desc := fmt.Sprintf("This is a product ordered by %s %s", selectedOrder.FirstName, selectedOrder.LastName)
	qrCodeParam := new(payment_request.QRCodeParameters)

	qrisCode := payment_request.QRCODECHANNELCODE_DANA
	qrCodeParam.ChannelCode = *payment_request.NewNullableQRCodeChannelCode(&qrisCode)
	qrCodeParam.ChannelProperties = payment_request.NewQRCodeChannelProperties()

	custId := strconv.FormatUint(selectedOrder.UserId, 10)
	metadata := map[string]interface{}{
		"user_id":  selectedOrder.UserId,
		"order_id": selectedOrder.ID,
		"notes":    selectedOrder.Note,
	}

	paymentRequestParameters := &payment_request.PaymentRequestParameters{
		Amount:      &amountFloat64,
		Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR,
		Description: *payment_request.NewNullableString(&desc),
		PaymentMethod: &payment_request.PaymentMethodParameters{
			Type:        payment_request.PAYMENTMETHODTYPE_QR_CODE,
			Reusability: payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
			QrCode:      *payment_request.NewNullableQRCodeParameters(qrCodeParam),
		},
		Items:      *paymentRequestBasketItems,
		CustomerId: *payment_request.NewNullableString(&custId),
		Metadata:   metadata,
	}

	resp, _, resErr := c.XenditClient.PaymentRequestApi.CreatePaymentRequest(ctx.Context()).
		PaymentRequestParameters(*paymentRequestParameters).IdempotencyKey(selectedOrder.Invoice).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to create new xendit transaction : %+v", resErr.FullError())
		return nil, fiber.ErrInternalServerError
	}

	// setelah itu tangkap semua response
	newXenditTransaction := new(entity.XenditTransactions)
	newXenditTransaction.ID = resp.Id
	newXenditTransaction.OrderId = selectedOrder.ID
	newXenditTransaction.ReferenceId = resp.ReferenceId
	newXenditTransaction.Amount = *resp.Amount
	newXenditTransaction.Currency = resp.Currency.String()
	newXenditTransaction.PaymentMethod = resp.PaymentMethod.Type.String()
	newXenditTransaction.PaymentMethodId = resp.PaymentMethod.Id
	newXenditTransaction.ChannelCode = resp.PaymentMethod.QrCode.Get().ChannelCode.Get().String()
	newXenditTransaction.QrString = *resp.PaymentMethod.QrCode.Get().ChannelProperties.QrString
	newXenditTransaction.Status = string(resp.Status)
	newXenditTransaction.FailureCode = resp.GetFailureCode()
	getMetadata := resp.GetMetadata()
	if getMetadata != nil {
		jsonMetadata, err := json.Marshal(metadata)
		if err != nil {
			c.Log.Warnf("failed to parse to json metadata : %+v", resErr.FullError())
			return nil, fiber.ErrInternalServerError
		}
		newXenditTransaction.Metadata = jsonMetadata
	}
	newXenditTransaction.Description = resp.GetDescription()
	expiresAt := resp.PaymentMethod.QrCode.Get().ChannelProperties.ExpiresAt.Format(time.DateTime)
	newXenditTransaction.ExpiresAt = expiresAt
	parseCreatedAt, err := ParseToRFC3339(resp.Created)
	if err != nil {
		c.Log.Warnf("failed to parse created_at into UTC : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newXenditTransaction.Created_At = parseCreatedAt

	parseUpdatedAt, err := ParseToRFC3339(resp.Updated)
	if err != nil {
		c.Log.Warnf("failed to parse updated_at into UTC : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newXenditTransaction.Updated_At = parseUpdatedAt
	if err := c.XenditTransactionRepository.Create(tx, newXenditTransaction); err != nil {
		c.Log.Warnf("Failed to insert xendit transaction into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "An error occurred on the server. Please try again later!")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.XenditTransactionToResponse(*newXenditTransaction), nil
}

func ParseToRFC3339(TimeRFC3339Nano string) (string, error) {
	parse, err := time.Parse(time.RFC3339Nano, TimeRFC3339Nano)
	if err != nil {
		return "", err
	}
	timeAtUTC := parse.UTC()
	parseToRFC3339 := timeAtUTC.Format(time.DateTime)
	return parseToRFC3339, nil
}

func (c *XenditTransactionQRCodeUseCase) GetTransaction(ctx *fiber.Ctx, request *model.GetXenditQRCodeTransaction) (*model.XenditTransactionResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newXenditTransaction := new(entity.XenditTransactions)
	if err := c.XenditTransactionRepository.FirstXenditTransactionByOrderId(tx, newXenditTransaction, request.OrderId, "Order", "Order.OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	resp, _, resErr := c.XenditClient.PaymentRequestApi.GetPaymentRequestByID(ctx.Context(), newXenditTransaction.ID).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to find xendit transaction : %+v", resErr.FullError())
		return nil, fiber.ErrInternalServerError
	}

	if newXenditTransaction.Status != string(resp.Status) {
		// update status payment
		hasPaymentStatusUpdated := false
		newXenditTransaction.Status = string(resp.Status)
		parseUpdatedAt, err := time.Parse(time.RFC3339Nano, resp.Updated)
		if err != nil {
			c.Log.Warnf("failed to parse updated_at into UTC : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		newXenditTransaction.Updated_At = parseUpdatedAt.Format(time.RFC3339)
		updatePaymentStatus := map[string]interface{}{
			"status":     string(resp.Status),
			"updated_at": parseUpdatedAt.Format(time.DateTime),
		}
		xenditTransactionObj := new(entity.XenditTransactions)
		xenditTransactionObj.ID = newXenditTransaction.ID
		if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, xenditTransactionObj, updatePaymentStatus); err != nil {
			c.Log.Warnf("failed to update xendit transaction : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// update juga di orders
		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED {
			// paid
			newXenditTransaction.Order.PaymentStatus = 2
			hasPaymentStatusUpdated = true
		}

		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_FAILED {
			// not paid
			newXenditTransaction.Order.PaymentStatus = 0
			hasPaymentStatusUpdated = true
		}

		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_CANCELED {
			// cancelled
			newXenditTransaction.Order.PaymentStatus = -1
			hasPaymentStatusUpdated = true
		}

		if hasPaymentStatusUpdated {
			updatePaymentStatus = map[string]interface{}{
				"payment_status": newXenditTransaction.Order.PaymentStatus,
				"updated_at":     time.Now().Format(time.DateTime),
			}
			orderObj := new(entity.Order)
			orderObj.ID = newXenditTransaction.Order.ID
			if err := c.OrderRepository.UpdateCustomColumns(tx, orderObj, updatePaymentStatus); err != nil {
				c.Log.Warnf("failed to update order payment status : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.XenditTransactionToResponse(*newXenditTransaction), nil
}
