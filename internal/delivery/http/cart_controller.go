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
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}
	auth := middleware.GetCurrentUser(ctx)
	request.UserID = auth.ID

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to add product to cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.CartResponse]{
		Code:   201,
		Status: "product successfully added to cart",
		Data:   response,
	})
}

func (c *CartController) Update(ctx *fiber.Ctx) error {
	getId := ctx.Params("cartItemId")
	cartItemId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert order_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert order_id to integer : %+v", err))
	}

	request := new(model.UpdateCartRequest)
	request.CartItemID = uint64(cartItemId)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.UpdateQuantity(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to update product quantity in the cart: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartResponse]{
		Code:   200,
		Status: "cart item updated successfully",
		Data:   response,
	})
}

func (c *CartController) GetAllCurrent(ctx *fiber.Ctx) error {
	request := new(model.GetAllCartByCurrentUserRequest)
	auth := middleware.GetCurrentUser(ctx)
	request.UserID = auth.ID

	response, err := c.UseCase.GetAllByCurrentUser(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to get all products from cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartResponse]{
		Code:   200,
		Status: "success to get all product from cart",
		Data:   response,
	})
}

func (c *CartController) Delete(ctx *fiber.Ctx) error {
	getId := ctx.Params("cartItemId")
	cartItemId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert cart_item_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert cart_item_id to integer : %+v", err))
	}

	request := new(model.DeleteCartRequest)
	request.CartItemID = uint64(cartItemId)

	response, err := c.UseCase.DeleteItem(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to delete product from cart : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CartItemResponse]{
		Code:   200,
		Status: "product removed from cart successfully",
		Data:   response,
	})
}
