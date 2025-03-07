package http

import (
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CartController struct {
	Log     *logrus.Logger
	UseCase *usecase.CartUseCase
}

func NewCartController(useCase *usecase.CartUseCase, logger *logrus.Logger) *CartController {
	return &CartController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *CartController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateCartRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	auth := middleware.GetCurrentUser(ctx)
	request.UserID = auth.ID

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to add product to cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.CartResponse]{
		Code:   201,
		Status: "Product successfully added to cart",
		Data:   response,
	})
}

func (c *CartController) Update(ctx *fiber.Ctx) error {
	getId := ctx.Params("cartItemId")
	cartItemId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order id : %+v", err)
		return err
	}

	request := new(model.UpdateCartRequest)
	request.CartItemID = uint64(cartItemId)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.UpdateQuantity(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to update product quantity in the cart: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartItemResponse]{
		Code:   200,
		Status: "Cart item updated successfully",
		Data:   response,
	})
}

func (c *CartController) GetAllCurrent(ctx *fiber.Ctx) error {
	request := new(model.GetAllCartByCurrentUserRequest)
	auth := middleware.GetCurrentUser(ctx)
	request.UserID = auth.ID

	response, err := c.UseCase.GetAllByCurrentUser(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to get all products from cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartResponse]{
		Code:   200,
		Status: "Success to get all product from cart",
		Data:   response,
	})
}

func (c *CartController) Delete(ctx *fiber.Ctx) error {
	getId := ctx.Params("cartItemId")
	cartItemId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order id : %+v", err)
		return err
	}

	request := new(model.DeleteCartRequest)
	request.CartItemID = uint64(cartItemId)

	response, err := c.UseCase.DeleteItem(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete product from cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartItemResponse]{
		Code:   200,
		Status: "Product removed from cart successfully",
		Data:   response,
	})
}
