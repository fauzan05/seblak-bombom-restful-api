package usecase

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
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
	UserRepository              *repository.UserRepository
	WalletRepository            *repository.WalletRepository
	XenditPayoutRepository      *repository.XenditPayoutRepository
	PayoutRepository            *repository.PayoutRepository
}

func NewXenditCallbackUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditClient *xendit.APIClient, xenditPayoutRepository *repository.XenditPayoutRepository,
	userRepository *repository.UserRepository, walletRepository *repository.WalletRepository,
	payoutRepository *repository.PayoutRepository) *XenditCallbackUseCase {
	return &XenditCallbackUseCase{
		DB:                          db,
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		XenditTransactionRepository: xenditTransactionRepository,
		XenditPayoutRepository:      xenditPayoutRepository,
		XenditClient:                xenditClient,
		UserRepository:              userRepository,
		WalletRepository:            walletRepository,
		PayoutRepository:            payoutRepository,
	}
}

func (c *XenditCallbackUseCase) UpdateStatusPaymentRequestCallback(ctx *fiber.Ctx, request *model.XenditGetPaymentRequestCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newXenditTransaction := new(entity.XenditTransactions)
	count, err := c.XenditTransactionRepository.FindXenditTransaction(tx, newXenditTransaction, request.Data.PaymentMethod.ID)
	if err != nil {
		c.Log.Warnf("Failed to get xendit transaction from database : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get xendit transaction from database : %+v", err))
	}

	if count > 0 {
		// update datanya
		if newXenditTransaction.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt.Format(time.DateTime)
			status := request.Data.Status
			orderId := newXenditTransaction.OrderId
			updateXenditTransaction := map[string]any{
				"status":     status,
				"updated_at": updatedAt,
			}

			*newXenditTransaction = entity.XenditTransactions{
				ID: newXenditTransaction.ID,
			}

			if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, newXenditTransaction, updateXenditTransaction); err != nil {
				c.Log.Warnf("Failed to update xendit transaction status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update xendit transaction status into database : %+v", err))
			}

			var payment_status helper.PaymentStatus
			if status == string(payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED) {
				payment_status = helper.PAID_PAYMENT
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_CANCELED) {
				payment_status = helper.CANCELLED_PAYMENT
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_FAILED) {
				payment_status = helper.FAILED_PAYMENT
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_EXPIRED) {
				payment_status = helper.EXPIRED_PAYMENT
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_PENDING) {
				payment_status = helper.PENDING_PAYMENT
			}

			updateOrderStatus := map[string]any{
				"payment_status": payment_status,
				"updated_at":     updatedAt,
			}

			newOrder := new(entity.Order)
			newOrder.ID = orderId
			if err := c.OrderRepository.UpdateCustomColumns(tx, newOrder, updateOrderStatus); err != nil {
				c.Log.Warnf("Failed to update order status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update order status into database : %+v", err))
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return nil
}

func (c *XenditCallbackUseCase) UpdateStatusPayoutRequestCallback(ctx *fiber.Ctx, request *model.XenditGetPayoutRequestCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newXenditPayout := new(entity.XenditPayout)
	newXenditPayout.ID = request.Data.PayoutId
	count, err := c.XenditPayoutRepository.FindFirstAndCount(tx, newXenditPayout)
	if err != nil {
		c.Log.Warnf("Failed to get xendit transaction from database : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get xendit transaction from database : %+v", err))
	}

	if count > 0 {
		// update datanya
		if newXenditPayout.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt.Format(time.DateTime)
			status := request.Data.Status
			updateXenditPayout := map[string]any{
				"status":     status,
				"updated_at": updatedAt,
			}

			*newXenditPayout = entity.XenditPayout{
				ID: newXenditPayout.ID,
			}

			if err := c.XenditPayoutRepository.UpdateCustomColumns(tx, newXenditPayout, updateXenditPayout); err != nil {
				c.Log.Warnf("Failed to update xendit payout status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update xendit payout status into database : %+v", err))
			}

			if request.Data.Status == "SUCCEEDED" {
				// update tb_payout
				newPayout := new(entity.Payout)
				if err := c.PayoutRepository.FindFirstPayoutByXenditPayoutId(tx, newPayout, newXenditPayout.ID); err != nil {
					c.Log.Warnf("Failed to get payout by xendit payout id from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get payout by xendit payout id from database : %+v", err))
				}

				if newPayout.ID < 1 {
					c.Log.Warnf("Payout not found!")
					return fiber.NewError(fiber.StatusNotFound, "Payout not found!")
				} else {
					// update payout
					updateStatus := map[string]any{
						"status": helper.PAYOUT_SUCCEEDED,
					}

					if err := c.PayoutRepository.UpdateCustomColumns(tx, newPayout, updateStatus); err != nil {
						c.Log.Warnf("Failed to update payout status in the database : %+v", err)
						return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update payout status in the database : %+v", err))
					}
				}
			}

			if request.Data.Status == "CANCELLED" || request.Data.Status == "FAILED" || request.Data.Status == "EXPIRED" || request.Data.Status == "REFUNDED" {
				// kembalikan saldonya
				newUser := new(entity.User)
				newUser.ID = newXenditPayout.UserID
				if err := c.UserRepository.FindWithPreloads(tx, newUser, "Wallet"); err != nil {
					c.Log.Warnf("Failed to find user wallet from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ailed to find user wallet from database : %+v", err))
				}

				resultBalance := newUser.Wallet.Balance + request.Data.Amount
				// update saldo
				updateBalance := map[string]any{
					"balance": resultBalance,
				}

				newWallet := new(entity.Wallet)
				newWallet.ID = newUser.Wallet.ID
				if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
					c.Log.Warnf("Failed to update wallet balance in the database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update wallet balance in the database : %+v", err))
				}

				// update tb_payout
				newPayout := new(entity.Payout)
				if err := c.PayoutRepository.FindFirstPayoutByXenditPayoutId(tx, newPayout, newXenditPayout.ID); err != nil {
					c.Log.Warnf("Failed to get payout by xendit payout id from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get payout by xendit payout id from database : %+v", err))
				}

				if newPayout.ID < 1 {
					c.Log.Warnf("Payout not found!")
					return fiber.NewError(fiber.StatusNotFound, "Payout not found!")
				} else {
					status := helper.PAYOUT_CANCELLED
					if request.Data.Status == "FAILED" {
						status = helper.PAYOUT_FAILED
					} else if request.Data.Status == "EXPIRED" {
						status = helper.PAYOUT_EXPIRED
					} else if request.Data.Status == "REFUNDED" {
						status = helper.PAYOUT_REFUNDED
					}

					// update payout
					updateStatus := map[string]any{
						"status": status,
					}

					if err := c.PayoutRepository.UpdateCustomColumns(tx, newPayout, updateStatus); err != nil {
						c.Log.Warnf("Failed to update payout status in the database : %+v", err)
						return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update payout status in the database : %+v", err))
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return nil
}
