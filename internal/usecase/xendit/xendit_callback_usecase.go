package usecase

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payment_request"
	"gorm.io/gorm"
)

type XenditCallbackUseCase struct {
	DB                          *gorm.DB
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	XenditClient                *xendit.APIClient
	OrderRepository             *repository.OrderRepository
	XenditTransactionRepository *repository.XenditTransctionRepository
}

func NewXenditCallbackUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditClient *xendit.APIClient) *XenditCallbackUseCase {
	return &XenditCallbackUseCase{
		DB:                          db,
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		XenditTransactionRepository: xenditTransactionRepository,
		XenditClient:                xenditClient,
	}
}

func (c *XenditCallbackUseCase) UpdateStatusPaymentRequestCallback(ctx *fiber.Ctx, request *model.XenditGetPaymentRequestCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	newXenditTransaction := new(entity.XenditTransactions)
	newXenditTransaction.ID = request.Data.PaymentRequestId
	count, err := c.XenditTransactionRepository.FindFirstAndCount(tx, newXenditTransaction)
	if err != nil {
		c.Log.Warnf("Failed to get xendit transaction from database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		// update datanya
		if newXenditTransaction.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt.Format(time.DateTime)
			status := request.Data.Status
			orderId := newXenditTransaction.OrderId
			updateXenditTransaction := map[string]interface{}{
				"status": status,
				"updated_at": updatedAt,
			}
			*newXenditTransaction = entity.XenditTransactions{
				ID: request.Data.PaymentRequestId,
			}
			if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, newXenditTransaction, updateXenditTransaction); err != nil {
				c.Log.Warnf("Failed to update xendit transaction status into database : %+v", err)
				return fiber.ErrInternalServerError
			}

			var payment_status int
			if status == string(payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED) {
				payment_status = 1
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_CANCELED) {
				payment_status = -1
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_FAILED) {
				payment_status = -3
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_EXPIRED) {
				payment_status = -2
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_PENDING) {
				payment_status = 0
			}

			updateOrderStatus := map[string]interface{}{
				"payment_status": payment_status,
				"updated_at": updatedAt,
			}

			newOrder := new(entity.Order)
			newOrder.ID = orderId
			if err := c.OrderRepository.UpdateCustomColumns(tx, newOrder, updateOrderStatus); err != nil {
				c.Log.Warnf("Failed to update order status into database : %+v", err)
				return fiber.ErrInternalServerError
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *XenditCallbackUseCase) UpdateStatusPaymentMethodCallback(ctx *fiber.Ctx, request *model.XenditGetPaymentMethodCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	newXenditTransaction := new(entity.XenditTransactions)
	count, err := c.XenditTransactionRepository.FindXenditTransactionByPaymentMethodId(tx, newXenditTransaction, request.Data.PaymentMethodId)
	if err != nil {
		c.Log.Warnf("Failed to get xendit transaction by payment method id from database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		// update datanya
		if newXenditTransaction.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt.Format(time.DateTime)
			status := request.Data.Status
			orderId := newXenditTransaction.OrderId
			xenditTransactionId := newXenditTransaction.ID
			updateXenditTransaction := map[string]interface{}{
				"status": status,
				"updated_at": updatedAt,
			}
			*newXenditTransaction = entity.XenditTransactions{
				ID: xenditTransactionId,
			}
			
			if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, newXenditTransaction, updateXenditTransaction); err != nil {
				c.Log.Warnf("Failed to update xendit transaction status into database : %+v", err)
				return fiber.ErrInternalServerError
			}

			var payment_status int
			if status == string(payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED) {
				payment_status = 1
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_CANCELED) {
				payment_status = -1
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_FAILED) {
				payment_status = -3
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_EXPIRED) {
				payment_status = -2
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_PENDING) {
				payment_status = 0
			}

			updateOrderStatus := map[string]interface{}{
				"payment_status": payment_status,
				"updated_at": updatedAt,
			}

			newOrder := new(entity.Order)
			newOrder.ID = orderId
			if err := c.OrderRepository.UpdateCustomColumns(tx, newOrder, updateOrderStatus); err != nil {
				c.Log.Warnf("Failed to update order status into database : %+v", err)
				return fiber.ErrInternalServerError
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}