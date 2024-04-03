package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DeliveryController struct {
	Log     *logrus.Logger
	UseCase *usecase.DeliveryUseCase
}

func NewDeliveryController(useCase *usecase.DeliveryUseCase, logger *logrus.Logger) *DeliveryController {
	return &DeliveryController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *DeliveryController) Create(ctx  *fiber.Ctx) error {
	request := new(model.CreateDeliveryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create/update delivery setting : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.DeliveryResponse]{
		Code:   201,
		Status: "Success to create a delivery settings",
		Data:   response,
	})
}

func (c *DeliveryController) Get(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Get(ctx.Context())
	if err != nil {
		c.Log.Warnf("Failed to get a delivery setting data : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DeliveryResponse]{
		Code:   200,
		Status: "Success to get a delivery settings",
		Data:   response,
	})
}

func (c *DeliveryController) Update(ctx  *fiber.Ctx) error {
	deliveryRequest := new(model.UpdateDeliveryRequest)
	if err := ctx.BodyParser(deliveryRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	getId := ctx.Params("deliveryId")
	deliveryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert category id : %+v", err)
		return err
	}
	deliveryRequest.ID = uint64(deliveryId)

	response, err := c.UseCase.Edit(ctx.Context(), deliveryRequest)
	if err != nil {
		c.Log.Warnf("Failed to update delivery setting : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DeliveryResponse]{
		Code:   200,
		Status: "Success to update a delivery settings",
		Data:   response,
	})
}

func (c *DeliveryController) Remove(ctx *fiber.Ctx) error {
	deliveryRequest := new(model.DeleteDeliveryRequest)
	getId := ctx.Params("deliveryId")
	deliveryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert category id : %+v", err)
		return err
	}
	deliveryRequest.ID = uint64(deliveryId)
	response, err := c.UseCase.Delete(ctx.Context(), deliveryRequest)
	if err != nil {
		c.Log.Warnf("Failed to remove a delivery setting by id : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to remove a delivery settings",
		Data:   response,
	})
}