package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletUseCase struct {
	DB                              *gorm.DB
	Log                             *logrus.Logger
	Validate                        *validator.Validate
	UserRepository                  *repository.UserRepository
	WalletRepository                *repository.WalletRepository
	WalletWithdrawRequestRepository *repository.WalletWithdrawRequestRepository
}

func NewWalletUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, walletRepository *repository.WalletRepository,
	walletWithdrawRequestRepository *repository.WalletWithdrawRequestRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:                              db,
		Log:                             log,
		Validate:                        validate,
		UserRepository:                  userRepository,
		WalletRepository:                walletRepository,
		WalletWithdrawRequestRepository: walletWithdrawRequestRepository,
	}
}

func (c *WalletUseCase) WithdrawByCustRequest(ctx context.Context, request *model.WithdrawWalletRequest) (*model.WithdrawWalletResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newUser := new(entity.User)
	newUser.ID = request.UserId
	if err := c.UserRepository.FindWithPreloads(tx, newUser, "Wallet"); err != nil {
		c.Log.Warnf("failed to get user by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get user by id : %+v", err))
	}

	if request.Amount > newUser.Wallet.Balance {
		c.Log.Warnf("your balance is insufficient to perform this withdraw transaction!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "your balance is insufficient to perform this withdraw transaction!")
	}

	newWithdrawRequest := new(entity.WalletWithdrawRequests)
	newWithdrawRequest.UserId = request.UserId
	newWithdrawRequest.Amount = request.Amount
	newWithdrawRequest.Method = request.Method
	newWithdrawRequest.BankName = request.BankName
	newWithdrawRequest.BankAcountName = request.BankAccountName
	newWithdrawRequest.BankAcountNumber = request.BankAccountNumber
	newWithdrawRequest.Status = request.Status
	newWithdrawRequest.Note = request.Note
	if err := c.WalletWithdrawRequestRepository.Create(tx, newWithdrawRequest); err != nil {
		c.Log.Warnf("failed to create new wallet withdraw request : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create new wallet withdraw request : %+v", err))
	}

	newWallet := new(entity.Wallet)
	newWallet.ID = newUser.Wallet.ID
	newWallet.UserId = newUser.ID
	newWallet.Balance = newUser.Wallet.Balance - request.Amount
	newWallet.Status = newUser.Wallet.Status
	if err := c.WalletRepository.Update(tx, newWallet); err != nil {
		c.Log.Warnf("failed to update wallet balance : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance : %+v", err))
	}

	newSaveWalletTransaction := new(helper_others.SaveWalletTransactionRequest)
	newSaveWalletTransaction.DB = tx
	newSaveWalletTransaction.UserId = request.UserId
	newSaveWalletTransaction.OrderId = nil
	newSaveWalletTransaction.Amount = request.Amount
	newSaveWalletTransaction.FlowType = enum_state.WALLET_FLOW_TYPE_DEBIT
	newSaveWalletTransaction.TransactionType = enum_state.WALLET_TRANSACTION_TYPE_WITHDRAW
	newSaveWalletTransaction.PaymentMethod = enum_state.PAYMENT_METHOD_CASH
	newSaveWalletTransaction.Status = enum_state.WALLET_TRANSACTION_STATUS_COMPLETED
	newSaveWalletTransaction.ReferenceNumber = ""
	newSaveWalletTransaction.Note = fmt.Sprintf("Withdraw request for Rp %f", request.Amount)
	newSaveWalletTransaction.AdminNote = ""
	newSaveWalletTransaction.ProcessedAt = nil
	newSaveWalletTransaction.ProcessedBy = nil

	err = helper_others.SaveWalletTransaction(newSaveWalletTransaction)
	if err != nil {
		c.Log.Warnf("failed to save wallet transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
	}

	if err := c.WalletWithdrawRequestRepository.FindWithPreloads(tx, newWithdrawRequest, "User"); err != nil {
		c.Log.Warnf("failed to get wallet by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get wallet by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.WalletWithdrawToResponse(newWithdrawRequest), nil
}

func (c *WalletUseCase) WithdrawByAdminApproval(ctx context.Context, request *model.WithdrawWalletApprovalRequest) (*model.WithdrawWalletResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newWithdrawRequest := new(entity.WalletWithdrawRequests)
	newWithdrawRequest.ID = request.ID
	if err := c.WalletWithdrawRequestRepository.FindById(tx, newWithdrawRequest); err != nil {
		c.Log.Warnf("failed to find withdraw wallet request by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find withdraw wallet request by id : %+v", err))
	}
	transactionStatus := enum_state.WALLET_TRANSACTION_STATUS_COMPLETED

	if request.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING {
		c.Log.Warnf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status))
	}

	if request.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_REJECTED {
		if newWithdrawRequest.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_REJECTED || newWithdrawRequest.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_APPROVED {
			c.Log.Warnf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status))
		}
		transactionStatus = enum_state.WALLET_TRANSACTION_STATUS_FAILED
	}

	if request.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_APPROVED {
		if newWithdrawRequest.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_REJECTED || newWithdrawRequest.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_APPROVED {
			c.Log.Warnf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("can't pending an wallet withdraw that has been %s!", newWithdrawRequest.Status))
		}
	}

	now := time.Now()
	newWithdrawRequest.Status = request.Status
	newWithdrawRequest.RejectionNotes = request.RejectionNotes
	newWithdrawRequest.ProcessedAt = &now
	newWithdrawRequest.ProcessedBy = &request.CurrentAdminId
	if err := c.WalletWithdrawRequestRepository.Update(tx, newWithdrawRequest); err != nil {
		c.Log.Warnf("failed to update wallet withdraw request : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet withdraw request : %+v", err))
	}

	walletTransactionStatus := enum_state.WALLET_TRANSACTION_STATUS_COMPLETED
	if request.Status == enum_state.WALLET_WITHDRAW_REQUEST_STATUS_REJECTED {
		walletTransactionStatus = enum_state.WALLET_TRANSACTION_STATUS_FAILED
		newWallet := new(entity.Wallet)
		newWallet.UserId = newWithdrawRequest.UserId
		if err := c.WalletRepository.FindFirstWalletByUserId(tx, newWallet, newWithdrawRequest.UserId, string(enum_state.ACTIVE_WALLET)); err != nil {
			c.Log.Warnf("failed to get wallet by user id : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get wallet by user id : %+v", err))
		}

		newWallet.Balance = newWallet.Balance + newWithdrawRequest.Amount
		newWallet.Status = enum_state.ACTIVE_WALLET
		if err := c.WalletRepository.Update(tx, newWallet); err != nil {
			c.Log.Warnf("failed to update wallet balance : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance : %+v", err))
		}
	}

	newSaveWalletTransaction := new(helper_others.SaveWalletTransactionRequest)
	newSaveWalletTransaction.DB = tx
	newSaveWalletTransaction.UserId = newWithdrawRequest.UserId
	newSaveWalletTransaction.OrderId = nil
	newSaveWalletTransaction.Amount = newWithdrawRequest.Amount
	newSaveWalletTransaction.FlowType = enum_state.WALLET_FLOW_TYPE_DEBIT
	newSaveWalletTransaction.TransactionType = enum_state.WALLET_TRANSACTION_TYPE_WITHDRAW
	newSaveWalletTransaction.PaymentMethod = enum_state.PAYMENT_METHOD_CASH
	newSaveWalletTransaction.Status = transactionStatus
	newSaveWalletTransaction.ReferenceNumber = ""
	newSaveWalletTransaction.Note = newWithdrawRequest.Note
	newSaveWalletTransaction.AdminNote = request.RejectionNotes
	newSaveWalletTransaction.ProcessedAt = &now
	newSaveWalletTransaction.ProcessedBy = &request.CurrentAdminId
	newSaveWalletTransaction.Status = walletTransactionStatus

	err = helper_others.SaveWalletTransaction(newSaveWalletTransaction)
	if err != nil {
		c.Log.Warnf("failed to save wallet transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
	}

	if err := c.WalletWithdrawRequestRepository.FindWithPreloads(tx, newWithdrawRequest, "User"); err != nil {
		c.Log.Warnf("failed to get wallet by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get wallet by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.WalletWithdrawToResponse(newWithdrawRequest), nil
}
