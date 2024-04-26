package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MidtransSnapOrderController struct {
	Log     *logrus.Logger
	UseCase *usecase.MidtransSnapOrderUseCase
}

func NewMidtransController(useCase *usecase.MidtransSnapOrderUseCase, logger *logrus.Logger) *MidtransSnapOrderController {
	return &MidtransSnapOrderController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *MidtransSnapOrderController) CreateSnap(ctx *fiber.Ctx) error {
	snapRequest := new(model.CreateMidtransSnapOrderRequest)
	if err := ctx.BodyParser(snapRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), snapRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new midtrans snap order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.MidtransSnapOrderResponse]{
		Code:   201,
		Status: "Success to create a new midtrans snap order",
		Data:   response,
	})
}
