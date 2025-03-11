package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProductReviewController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductReviewUseCase
}

func NewProductReviewController(useCase *usecase.ProductReviewUseCase, logger *logrus.Logger) *ProductReviewController {
	return &ProductReviewController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductReviewController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetCurrentUser(ctx)
	request := new(model.CreateProductReviewRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}

	request.UserId = auth.ID
	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create a product review : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.ProductReviewResponse]{
		Code:   201,
		Status: "Success to create a product review",
		Data:   response,
	})
}
