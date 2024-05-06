package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ApplicationController struct {
	Log     *logrus.Logger
	UseCase *usecase.ApplicationUseCase
}

func NewApplicationController(useCase *usecase.ApplicationUseCase, logger *logrus.Logger) *ApplicationController {
	return &ApplicationController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ApplicationController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateApplicationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create new application : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   201,
		Status: "Success to create a new application",
		Data:   response,
	})
}

func (c *ApplicationController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateApplicationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Edit(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create new application : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   200,
		Status: "Success to update an application",
		Data:   response,
	})
}

func (c *ApplicationController) Get(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Get(ctx.Context())
	if err != nil {
		c.Log.Warnf("Failed to get an application : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   200,
		Status: "Success to get an application",
		Data:   response,
	})
}