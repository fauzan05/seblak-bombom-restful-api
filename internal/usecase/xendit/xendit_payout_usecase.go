package usecase

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payout"
	"gorm.io/gorm"
)

type XenditPayoutUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	XenditClient           *xendit.APIClient
	XenditPayoutRepository *repository.XenditPayoutRepository
	WalletRepository       *repository.WalletRepository
	UserRepository         *repository.UserRepository
}

func NewXenditPayoutUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	xenditPayoutRepository *repository.XenditPayoutRepository,
	xenditClient *xendit.APIClient, walletRepository *repository.WalletRepository,
	userRepository *repository.UserRepository) *XenditPayoutUseCase {
	return &XenditPayoutUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		XenditPayoutRepository: xenditPayoutRepository,
		WalletRepository:       walletRepository,
		UserRepository:         userRepository,
		XenditClient:           xenditClient,
	}
}

func (c *XenditPayoutUseCase) AddPayout(ctx *fiber.Ctx, request *model.CreateXenditPayout, tx *gorm.DB) (*model.XenditPayoutResponse, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	// cek apakah saldo melebihi yang ada
	newUser := new(entity.User)
	newUser.ID = request.UserId
	if err := c.UserRepository.FindWithPreloads(tx, newUser, "Wallet"); err != nil {
		c.Log.Warnf("failed to find wallet by user id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find wallet by user id : %+v", err))
	}

	if request.Amount > newUser.Wallet.Balance {
		c.Log.Warnf("your balance is insufficient to perform this transaction!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "your balance is insufficient to perform this transaction!")
	}

	resultBalance := newUser.Wallet.Balance - request.Amount

	// update saldo
	updateBalance := map[string]any{
		"balance": resultBalance,
	}

	newWallet := new(entity.Wallet)
	newWallet.ID = newUser.Wallet.ID
	if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
		c.Log.Warnf("failed to update wallet balance in the database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance in the database : %+v", err))
	}

	milli := time.Now().UnixMilli()
	idempotencyKey := fmt.Sprintf("disb-%d", milli)
	referenceId := fmt.Sprintf("payout-%d-%d", request.UserId, milli)
	channelProperties := payout.NewDigitalPayoutChannelProperties(request.AccountNumber)
	accountHolderName := payout.NewNullableString(&request.AccountHolderName)
	channelProperties.AccountHolderName = *accountHolderName
	createPayoutRequest := payout.NewCreatePayoutRequest(referenceId, request.ChannelCode, *channelProperties, request.Amount, request.Currency)
	receiptNotification := payout.NewReceiptNotification()
	emailSlice := []string{newUser.Email}
	receiptNotification.EmailTo = emailSlice
	createPayoutRequest.Description = &request.Description
	resp, _, resErr := c.XenditClient.PayoutApi.CreatePayout(ctx.Context()).
		IdempotencyKey(idempotencyKey).
		CreatePayoutRequest(*createPayoutRequest).Execute()

	if resErr != nil {
		c.Log.Warnf("failed to create new xendit payout : %+v", resErr.FullError())
		return nil, fiber.NewError(helper.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to create new xendit payout : %+v", resErr.FullError()))
	}

	newXenditPayout := new(entity.XenditPayout)
	newXenditPayout.ID = resp.Payout.Id
	newXenditPayout.UserID = request.UserId
	newXenditPayout.BusinessID = resp.Payout.GetBusinessId()
	newXenditPayout.ReferenceID = resp.Payout.GetReferenceId()
	newXenditPayout.Amount = resp.Payout.GetAmount()
	newXenditPayout.Currency = resp.Payout.GetCurrency()
	newXenditPayout.Description = resp.Payout.GetDescription()
	newXenditPayout.ChannelCode = resp.Payout.GetChannelCode()
	newXenditPayout.AccountNumber = resp.Payout.ChannelProperties.GetAccountNumber()
	newXenditPayout.AccountHolderName = resp.Payout.ChannelProperties.GetAccountHolderName()
	newXenditPayout.Status = resp.Payout.GetStatus()
	estimatedArrivalTime := resp.Payout.GetEstimatedArrivalTime()
	newXenditPayout.EstimatedArrival = estimatedArrivalTime

	if err := c.XenditPayoutRepository.Create(tx, newXenditPayout); err != nil {
		c.Log.Warnf("failed to insert xendit payout into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to insert xendit payout into database : %+v", err))
	}

	return converter.XenditPayoutToResponse(newXenditPayout), nil
}

func (c *XenditPayoutUseCase) GetBalance(ctx *fiber.Ctx) (*model.GetWithdrawableBalanceResponse, error) {
	tx := c.DB.WithContext(ctx.Context())
	var balance *float32
	resp, _, resErr := c.XenditClient.BalanceApi.GetBalance(ctx.Context()).Execute()
	if resErr != nil {
		c.Log.Warnf("failed to get xendit balance : %+v", resErr.FullError())
		return nil, fiber.NewError(helper.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to get xendit balance : %+v", resErr.FullError()))
	}

	newWallet := new(entity.Wallet)
	getActiveBalance, err := c.WalletRepository.FindAllActiveBalance(tx, newWallet)
	if err != nil {
		c.Log.Warnf("failed to count balance on active wallet : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to count balance on active wallet : %+v", err))
	}

	respBalance := resp.GetBalance()
	pointerBalance := &respBalance
	result := *pointerBalance - *getActiveBalance
	balance = &result

	return converter.WithdrawableBalanceResponse(balance, getActiveBalance), nil
}

func (c *XenditPayoutUseCase) GetPayoutById(ctx *fiber.Ctx, request *model.GetPayoutById) (*model.XenditPayoutResponse, error) {
	tx := c.DB.WithContext(ctx.Context())

	resp, _, resErr := c.XenditClient.PayoutApi.GetPayoutById(ctx.Context(), request.PayoutId).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to get xendit payout by id : %+v", resErr.FullError())
		return nil, fiber.NewError(helper.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to get xendit payout by id : %+v", resErr.FullError()))
	}

	newXenditPayout := new(entity.XenditPayout)
	newXenditPayout.ID = resp.Payout.Id
	if err := c.XenditPayoutRepository.FindWith2Preloads(tx, newXenditPayout, "User", "User.Wallet"); err != nil {
		c.Log.Warnf("failed to get xendit payout from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get xendit payout from database : %+v", err))
	}
	newXenditPayout.BusinessID = resp.Payout.GetBusinessId()
	newXenditPayout.ReferenceID = resp.Payout.GetReferenceId()
	newXenditPayout.Amount = resp.Payout.GetAmount()
	newXenditPayout.Currency = resp.Payout.GetCurrency()
	newXenditPayout.Description = resp.Payout.GetDescription()
	newXenditPayout.ChannelCode = resp.Payout.GetChannelCode()
	newXenditPayout.AccountNumber = resp.Payout.ChannelProperties.GetAccountNumber()
	newXenditPayout.AccountHolderName = resp.Payout.ChannelProperties.GetAccountHolderName()
	newXenditPayout.Status = resp.Payout.GetStatus()
	estimatedArrivalTime := resp.Payout.GetEstimatedArrivalTime()
	newXenditPayout.EstimatedArrival = estimatedArrivalTime

	return converter.XenditPayoutToResponse(newXenditPayout), nil
}

func (c *XenditPayoutUseCase) CancelPayout(ctx *fiber.Ctx, request *model.CancelXenditPayout) (*model.XenditPayoutResponse, error) {
	tx := c.DB.WithContext(ctx.Context())

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	resp, _, resErr := c.XenditClient.PayoutApi.CancelPayout(ctx.Context(), request.PayoutId).
		Execute()

	if resErr != nil {
		c.Log.Warnf("failed to cancel xendit payout : %+v", resErr.FullError())
		return nil, fiber.NewError(helper.SetFiberStatusCode(resErr.Status()), fmt.Sprintf("failed to cancel xendit payout : %+v", resErr.FullError()))
	}

	newXenditPayout := new(entity.XenditPayout)
	newXenditPayout.ID = resp.Payout.Id
	if err := c.XenditPayoutRepository.FindFirst(tx, newXenditPayout); err != nil {
		c.Log.Warnf("failed to get xendit payout from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get xendit payout from database : %+v", err))
	}
	newXenditPayout.BusinessID = resp.Payout.GetBusinessId()
	newXenditPayout.ReferenceID = resp.Payout.GetReferenceId()
	newXenditPayout.Amount = resp.Payout.GetAmount()
	newXenditPayout.Currency = resp.Payout.GetCurrency()
	newXenditPayout.Description = resp.Payout.GetDescription()
	newXenditPayout.ChannelCode = resp.Payout.GetChannelCode()
	newXenditPayout.AccountNumber = resp.Payout.ChannelProperties.GetAccountNumber()
	newXenditPayout.AccountHolderName = resp.Payout.ChannelProperties.GetAccountHolderName()
	newXenditPayout.Status = resp.Payout.GetStatus()
	estimatedArrivalTime := resp.Payout.GetEstimatedArrivalTime()
	newXenditPayout.EstimatedArrival = estimatedArrivalTime

	return converter.XenditPayoutToResponse(newXenditPayout), nil
}
