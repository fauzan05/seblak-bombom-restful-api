package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MidtransController struct {
	Log *logrus.Logger
	UseCase *usecase.MidtransUseCase
}

func NewMidtransController(useCase *usecase.MidtransUseCase, logger *logrus.Logger) *MidtransController {
	return &MidtransController{
		Log: logger,
		UseCase: useCase,
	}
}

func (c *MidtransController) CreateSnap(ctx *fiber.Ctx) error {
	snapRequest := new(model.CreateSnapRequest)
	if err := ctx.BodyParser(snapRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), snapRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new midtrans snap : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.SnapResponse]{
		Code:   201,
		Status: "Success to create a new snap",
		Data:   response,
	})
}