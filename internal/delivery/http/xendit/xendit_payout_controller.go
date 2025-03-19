package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase/xendit"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type XenditPayoutController struct {
	Log     *logrus.Logger
	UseCase *usecase.XenditPayoutUseCase
}

func NewXenditPayoutController(useCase *usecase.XenditPayoutUseCase, logger *logrus.Logger) *XenditPayoutController {
	return &XenditPayoutController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *XenditPayoutController) Create(ctx *fiber.Ctx) error {
	getId := ctx.Params("userId")
	userId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert user_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert user_id to integer : %+v", err))
	}

	xenditPayoutRequest := new(model.CreateXenditPayout)
	xenditPayoutRequest.UserId = uint64(userId)
	if err := ctx.BodyParser(xenditPayoutRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}

	response, err := c.UseCase.AddPayout(ctx, xenditPayoutRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new xendit payout : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.XenditPayoutResponse]{
		Code:   201,
		Status: "Success to create a new xendit payout",
		Data:   response,
	})
}

func (c *XenditPayoutController) GetAdminBalance(ctx *fiber.Ctx) error {
	balance, err := c.UseCase.GetBalance(ctx)
	if err != nil {
		c.Log.Warnf("Failed to get withdrawable balance : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.GetWithdrawableBalanceResponse]{
		Code:   200,
		Status: "Success to get withdrawable balance",
		Data:   balance,
	})
}

func (c *XenditPayoutController) GetPayoutById(ctx *fiber.Ctx) error {
	payoutId := ctx.Params("payoutId")
	xenditPayoutRequest := new(model.GetPayoutById)
	xenditPayoutRequest.PayoutId = payoutId
	balance, err := c.UseCase.GetPayoutById(ctx, xenditPayoutRequest)
	if err != nil {
		c.Log.Warnf("Failed to get payout by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.XenditPayoutResponse]{
		Code:   200,
		Status: "Success to get payout by id",
		Data:   balance,
	})
}

func (c *XenditPayoutController) Cancel(ctx *fiber.Ctx) error {
	payoutId := ctx.Params("payoutId")
	
	xenditPayoutRequest := new(model.CancelXenditPayout)
	xenditPayoutRequest.PayoutId = payoutId

	response, err := c.UseCase.CancelPayout(ctx, xenditPayoutRequest)
	if err != nil {
		c.Log.Warnf("Failed to cancel xendit payout by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.XenditPayoutResponse]{
		Code:   200,
		Status: "Success to cancel xendit payout by id",
		Data:   response,
	})
}