package usecase

import (
	"encoding/json"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
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

func (c *XenditTransactionQRCodeUseCase) Add(ctx *fiber.Ctx, request *model.CreateXenditTransaction, tx *gorm.DB) (*model.XenditTransactionResponse, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	// get order id
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("failed to find order by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find order by id : %+v", err))
	}

	paymentRequestBasketItems := new([]payment_request.PaymentRequestBasketItem)
	for _, product := range selectedOrder.OrderProducts {
		refId := strconv.FormatUint(product.ProductId, 10)
		itemType := string(enum_state.ITEM_TYPE_PHYSICAL_PRODUCT)
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
		itemType := string(enum_state.ITEM_TYPE_DELIVERY_FEE)
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
		itemType := string(enum_state.ITEM_TYPE_DISCOUNT)
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

	amountFloat64 := float64(selectedOrder.TotalFinalPrice)
	desc := fmt.Sprintf("This is a product ordered by %s %s", selectedOrder.FirstName, selectedOrder.LastName)
	qrCodeParam := new(payment_request.QRCodeParameters)

	qrisCode := payment_request.QRCODECHANNELCODE_DANA
	qrCodeParam.ChannelCode = *payment_request.NewNullableQRCodeChannelCode(&qrisCode)
	qrCodeParam.ChannelProperties = payment_request.NewQRCodeChannelProperties()
	setExpiresAt := time.Now().Add(5 * time.Minute)
	qrCodeParam.ChannelProperties.ExpiresAt = &setExpiresAt

	custId := strconv.FormatUint(selectedOrder.UserId, 10)
	metadata := map[string]any{
		"user_id":   selectedOrder.UserId,
		"order_id":  selectedOrder.ID,
		"notes":     selectedOrder.Note,
		"time_zone": request.TimeZone.String(),
		"lang":      request.Lang,
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

	idempotencyKey := fmt.Sprintf("%d-%s", selectedOrder.ID, selectedOrder.Invoice)
	resp, _, resErr := c.XenditClient.PaymentRequestApi.CreatePaymentRequest(ctx.Context()).
		PaymentRequestParameters(*paymentRequestParameters).IdempotencyKey(idempotencyKey).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to create new xendit transaction : %+v", resErr.FullError())
		return nil, fiber.NewError(helper_others.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to create new xendit transaction : %+v", resErr.FullError()))
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
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse to json metadata : %+v", resErr.FullError()))
		}
		newXenditTransaction.Metadata = jsonMetadata
	}

	newXenditTransaction.Description = resp.GetDescription()
	expiresAt := resp.PaymentMethod.QrCode.Get().ChannelProperties.ExpiresAt
	newXenditTransaction.ExpiresAt = *expiresAt
	parseCreatedAt, err := ParseToRFC3339(resp.Created)
	if err != nil {
		c.Log.Warnf("failed to parse created_at into UTC : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse created_at into UTC : %+v", err))
	}

	newXenditTransaction.CreatedAt = *parseCreatedAt
	parseUpdatedAt, err := ParseToRFC3339(resp.Updated)
	if err != nil {
		c.Log.Warnf("failed to parse updated_at into UTC : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse updated_at into UTC : %+v", err))
	}

	newXenditTransaction.UpdatedAt = *parseUpdatedAt
	if err := c.XenditTransactionRepository.Create(tx, newXenditTransaction); err != nil {
		c.Log.Warnf("failed to insert xendit transaction into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "An error occurred on the server. Please try again later!")
	}

	return converter.XenditTransactionToResponse(*newXenditTransaction), nil
}

func ParseToRFC3339(TimeRFC3339Nano string) (*time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, TimeRFC3339Nano)
	if err != nil {
		return nil, err
	}
	utcTime := parsedTime.UTC()
	return &utcTime, nil
}

func (c *XenditTransactionQRCodeUseCase) GetTransaction(ctx *fiber.Ctx, request *model.GetXenditQRCodeTransaction) (*model.XenditTransactionResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newXenditTransaction := new(entity.XenditTransactions)
	if err := c.XenditTransactionRepository.FirstXenditTransactionByOrderId(tx, newXenditTransaction, request.OrderId, "Order", "Order.OrderProducts"); err != nil {
		c.Log.Warnf("failed to find order by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find order by id : %+v", err))
	}

	resp, _, resErr := c.XenditClient.PaymentRequestApi.GetPaymentRequestByID(ctx.Context(), newXenditTransaction.ID).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to find xendit transaction : %+v", resErr.FullError())
		return nil, fiber.NewError(helper_others.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to find xendit transaction : %+v", resErr.FullError()))
	}

	if newXenditTransaction.Status != string(resp.Status) && newXenditTransaction.Status != string(payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED) {
		// update status payment
		hasPaymentStatusUpdated := false
		newXenditTransaction.Status = string(resp.Status)
		parseUpdatedAt, err := time.Parse(time.RFC3339Nano, resp.Updated)
		if err != nil {
			c.Log.Warnf("failed to parse updated_at into UTC : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse updated_at into UTC : %+v", err))
		}

		order_status := ""

		newXenditTransaction.UpdatedAt = parseUpdatedAt
		updatePaymentStatus := map[string]any{
			"status":     string(resp.Status),
			"updated_at": parseUpdatedAt.Format(time.DateTime),
		}

		if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, newXenditTransaction, updatePaymentStatus); err != nil {
			c.Log.Warnf("failed to update xendit transaction : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update xendit transaction : %+v", err))
		}

		// update juga di orders
		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED {
			// paid
			newXenditTransaction.Order.PaymentStatus = enum_state.PAID_PAYMENT
			hasPaymentStatusUpdated = true
		}

		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_FAILED {
			// not paid
			newXenditTransaction.Order.PaymentStatus = enum_state.FAILED_PAYMENT
			hasPaymentStatusUpdated = true
			order_status = string(enum_state.ORDER_CANCELLED)
		}

		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_CANCELED {
			// cancelled
			newXenditTransaction.Order.PaymentStatus = enum_state.CANCELLED_PAYMENT
			hasPaymentStatusUpdated = true
			order_status = string(enum_state.ORDER_CANCELLED)
		}

		if resp.Status == payment_request.PAYMENTREQUESTSTATUS_EXPIRED {
			// expired
			newXenditTransaction.Order.PaymentStatus = enum_state.EXPIRED_PAYMENT
			hasPaymentStatusUpdated = true
			order_status = string(enum_state.ORDER_CANCELLED)
		}

		if hasPaymentStatusUpdated {
			updatePaymentStatus = map[string]any{
				"payment_status": newXenditTransaction.Order.PaymentStatus,
				"updated_at":     time.Now().Format(time.DateTime),
			}

			if order_status != "" {
				updatePaymentStatus["order_status"] = order_status
			}

			orderObj := new(entity.Order)
			orderObj.ID = newXenditTransaction.Order.ID
			if err := c.OrderRepository.UpdateCustomColumns(tx, orderObj, updatePaymentStatus); err != nil {
				c.Log.Warnf("failed to update order payment status : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update order payment status : %+v", err))
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.XenditTransactionToResponse(*newXenditTransaction), nil
}
