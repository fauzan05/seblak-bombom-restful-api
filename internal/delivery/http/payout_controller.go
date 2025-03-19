package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PayoutController struct {
	Log     *logrus.Logger
	UseCase *usecase.PayoutUseCase
}

func NewPayoutController(useCase *usecase.PayoutUseCase, logger *logrus.Logger) *PayoutController {
	return &PayoutController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *PayoutController) Create(ctx *fiber.Ctx) error {
	getId := ctx.Params("userId")
	userId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert user_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert user_id to integer : %+v", err))
	}

	request := new(model.CreatePayoutRequest)
	request.UserId = uint64(userId)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create a new payout request : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.PayoutResponse]{
		Code:   201,
		Status: "Success to create a new payout request",
		Data:   response,
	})
}
