package usecase

import (
	"database/sql"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	xenditUseCase "seblak-bombom-restful-api/internal/usecase/xendit"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayoutUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	Validate            *validator.Validate
	PayoutRepository    *repository.PayoutRepository
	XenditPayoutUseCase *xenditUseCase.XenditPayoutUseCase
	WalletRepository    *repository.WalletRepository
	UserRepository      *repository.UserRepository
}

func NewPayoutUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	payoutRepository *repository.PayoutRepository, xenditPayoutUseCase *xenditUseCase.XenditPayoutUseCase,
	walletRepository *repository.WalletRepository, userRepository *repository.UserRepository) *PayoutUseCase {
	return &PayoutUseCase{
		DB:                  db,
		Log:                 log,
		Validate:            validate,
		PayoutRepository:    payoutRepository,
		XenditPayoutUseCase: xenditPayoutUseCase,
		UserRepository:      userRepository,
		WalletRepository:    walletRepository,
	}
}

func (c *PayoutUseCase) Add(ctx *fiber.Ctx, request *model.CreatePayoutRequest) (*model.PayoutResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	var xenditPayoutResponse *model.XenditPayoutResponse
	var xenditPayoutId sql.NullString
	var status helper.PayoutStatus
	if request.Method == helper.PAYOUT_METHOD_ONLINE {
		result, err := c.XenditPayoutUseCase.AddPayout(ctx, request.XenditPayoutRequest, tx)
		if err != nil {
			return nil, err
		}
		xenditPayoutResponse = result
		xenditPayoutId = sql.NullString{String: result.ID, Valid: true}
		status = helper.PayoutStatus(result.Status)
	} else {
		status = helper.PAYOUT_ACCEPTED
		newWallet := new(entity.Wallet)
		count, err := c.WalletRepository.FindAndCountFirstWalletByUserId(tx, newWallet, request.UserId, "active")
		if err != nil {
			c.Log.Warnf("Failed to find wallet by user id : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find wallet by user id : %+v", err))
		}

		if count < 1 {
			c.Log.Warnf("The selected wallet is not found!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "The selected wallet is not found!")
		}

		if request.Amount > newWallet.Balance {
			c.Log.Warnf("Your balance is insufficient to perform this transaction!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Your balance is insufficient to perform this transaction!")
		}

		resultBalance := newWallet.Balance - request.Amount
		updateBalance := map[string]any{
			"balance": resultBalance,
		}

		if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
			c.Log.Warnf("Failed to update wallet balance in the database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update wallet balance in the database : %+v", err))
		}
	}

	newPayout := new(entity.Payout)
	newPayout.UserId = request.UserId
	newPayout.XenditPayoutId = xenditPayoutId
	newPayout.Amount = request.Amount
	newPayout.Currency = request.Currency
	newPayout.Method = request.Method
	newPayout.Status = status
	newPayout.Notes = request.Notes

	if err := c.PayoutRepository.Create(tx, newPayout); err != nil {
		c.Log.Warnf("Failed to create payout request into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create payout request into database : %+v", err))
	}

	if xenditPayoutResponse != nil {
		if newPayout.XenditPayout == nil {
			newPayout.XenditPayout = &entity.XenditPayout{}
		}
		helper.PrintStructFields(xenditPayoutResponse)
		newPayout.XenditPayout.ID = xenditPayoutResponse.ID
		newPayout.XenditPayout.UserID = xenditPayoutResponse.UserId
		newPayout.XenditPayout.BusinessID = xenditPayoutResponse.BusinessId
		newPayout.XenditPayout.ReferenceID = xenditPayoutResponse.ReferenceId
		newPayout.XenditPayout.Amount = xenditPayoutResponse.Amount
		newPayout.XenditPayout.Currency = xenditPayoutResponse.Currency
		newPayout.XenditPayout.Description = xenditPayoutResponse.Description
		newPayout.XenditPayout.ChannelCode = xenditPayoutResponse.ChannelCode
		newPayout.XenditPayout.AccountNumber = xenditPayoutResponse.AccountNumber
		newPayout.XenditPayout.AccountHolderName = xenditPayoutResponse.AccountHolderName
		newPayout.XenditPayout.Status = xenditPayoutResponse.Status
		newPayout.XenditPayout.CreatedAt = xenditPayoutResponse.CreatedAt
		newPayout.XenditPayout.UpdatedAt = xenditPayoutResponse.UpdatedAt
		newPayout.XenditPayout.EstimatedArrival = xenditPayoutResponse.EstimatedArrival
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.PayoutToResponse(newPayout), nil
}
