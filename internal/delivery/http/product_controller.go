package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductUseCase
}

func NewProductController(useCase *usecase.ProductUseCase, logger *logrus.Logger) *ProductController {
	return &ProductController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateProductRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create new product : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   201,
		Status: "Success to create a new product",
		Data:   response,
	})
}

func (c *ProductController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("productId")
	productId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert product id : %+v", err)
		return err
	}
	productRequest := new(model.GetProductRequest)
	productRequest.ID = uint64(productId)

	response, err := c.UseCase.Get(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to find product by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   200,
		Status: "Success to get an product",
		Data:   response,
	})
}

func (c *ProductController) GetAll(ctx *fiber.Ctx) error {
	response, err := c.UseCase.GetAll(ctx.Context())
	if err != nil {
		c.Log.Warnf("Failed to find all products : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.ProductResponse]{
		Code:   200,
		Status: "Success to get all products",
		Data:   response,
	})
}

func (c *ProductController) Edit(ctx *fiber.Ctx) error {
	getId := ctx.Params("productId")
	productId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert product id : %+v", err)
		return err
	}

	productRequest := new(model.UpdateProductRequest)
	if err := ctx.BodyParser(productRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	productRequest.ID = uint64(productId)

	response, err := c.UseCase.Update(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to update product by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   200,
		Status: "Success to update an product by id",
		Data:   response,
	})
}

func (c *ProductController) Remove(ctx *fiber.Ctx) error {
	getId := ctx.Params("productId")
	productId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert product id : %+v", err)
		return err
	}
	productRequest := new(model.DeleteProductRequest)
	productRequest.ID = uint64(productId)

	response, err := c.UseCase.Delete(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to delete product by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete an product by id",
		Data:   response,
	})
}
