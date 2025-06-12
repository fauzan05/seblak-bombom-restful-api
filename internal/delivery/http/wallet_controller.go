package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

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
		c.Log.Warnf("failed to withdraw wallet request : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.WithdrawWalletResponse]{
		Code:   201,
		Status: "success to withdraw wallet request",
		Data:   response,
	})
}

func (c *WalletController) WithdrawAdminApproval(ctx *fiber.Ctx) error {
	request := new(model.WithdrawWalletApprovalRequest)
	getId := ctx.Params("withdrawRequestId")
	withdrawRequestId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert withdraw_request_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert withdraw_request_id to integer : %+v", err))
	}
	request.ID = uint64(withdrawRequestId)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	auth := middleware.GetCurrentUser(ctx)
	request.CurrentAdminId = auth.ID
	response, err := c.UseCase.WithdrawByAdminApproval(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to approval withdraw wallet balance : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.WithdrawWalletResponse]{
		Code:   200,
		Status: "success to approval withdraw wallet balance",
		Data:   response,
	})
}