package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	Log     *logrus.Logger
	UseCase *usecase.WalletUseCase
}

func NewWalletController(useCase *usecase.WalletUseCase, logger *logrus.Logger) *WalletController {
	return &WalletController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *WalletController) WithdrawCustRequest(ctx *fiber.Ctx) error {
	request := new(model.WithdrawWalletRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.WithdrawByCustRequest(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to withdraw wallet balance : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.WalletResponse]{
		Code:   200,
		Status: "success to withdraw wallet balance",
		Data:   response,
	})
}

func (c *WalletController) WithdrawCustApproval(ctx *fiber.Ctx) error {
	request := new(model.WithdrawWalletApprovalRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.WithdrawByCustApproval(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to withdraw wallet balance : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.WalletResponse]{
		Code:   200,
		Status: "success to withdraw wallet balance",
		Data:   response,
	})
}