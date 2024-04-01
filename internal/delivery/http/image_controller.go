package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ImageController struct {
	Log     *logrus.Logger
	UseCase *usecase.ImageUseCase
}

func NewImageController(useCase *usecase.ImageUseCase, logger *logrus.Logger) *ImageController {
	return &ImageController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ImageController) Creates(ctx *fiber.Ctx) error {
	request := new(model.AddImagesRequest)
	if err := ctx.BodyParser(&request.Images); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to add images : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*[]model.ImageResponse]{
		Code:   201,
		Status: "Success to create a new category",
		Data:   response,
	})
}