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

type OrderController struct {
	Log     *logrus.Logger
	UseCase *usecase.OrderUseCase
}

func NewOrderController(useCase *usecase.OrderUseCase, logger *logrus.Logger) *OrderController {
	return &OrderController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *OrderController) Create(ctx *fiber.Ctx) error {
	orderRequest := new(model.CreateOrderRequest)
	if err := ctx.BodyParser(orderRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}
	
	auth := middleware.GetCurrentUser(ctx)
	orderRequest.UserId = auth.ID
	orderRequest.FirstName = auth.FirstName
	orderRequest.LastName = auth.LastName
	orderRequest.Email = auth.Email
	orderRequest.Phone = auth.Phone
	if auth.Addresses == nil {
		c.Log.Warnf("Address not found/selected!")
		return fiber.NewError(fiber.StatusBadRequest, "Address not found/selected!")
	}

	for _, address := range auth.Addresses {
		if address.IsMain {
			orderRequest.CompleteAddress = address.CompleteAddress
			orderRequest.DeliveryId = address.Delivery.ID
		}
	}

	orderRequest.CurrentBalance = auth.Wallet.Balance
	response, err := c.UseCase.Add(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("Failed to create a new order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   201,
		Status: "Success to create a new order",
		Data:   response,
	})
}

func (c *OrderController) GetAllCurrent(ctx *fiber.Ctx) error {
	auth := middleware.GetCurrentUser(ctx)
	orderRequest := new(model.GetOrderByCurrentRequest)
	orderRequest.ID = auth.ID
	response, err := c.UseCase.GetAllCurrent(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("Failed to get all orders by current user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.OrderResponse]{
		Code:   200,
		Status: "Success to get all orders by current user",
		Data:   response,
	})
}

func (c *OrderController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	getId := ctx.Params("orderId")
	orderId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert order_id to integer : %+v", err))
	}

	orderRequest := new(model.UpdateOrderRequest)
	orderRequest.ID = uint64(orderId)
	if err := ctx.BodyParser(orderRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}
	
	response, err := c.UseCase.EditOrderStatus(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("Failed to update order status by selected order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   200,
		Status: "Success to update order status by selected order",
		Data:   response,
	})
}

func (c *OrderController) GetAllByUserId(ctx *fiber.Ctx) error {
	getId := ctx.Params("userId")
	userId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order id : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert order id : %+v", err))
	}

	orderRequest := new(model.GetOrdersByUserIdRequest)
	orderRequest.ID = uint64(userId)
	response, err := c.UseCase.GetByUserId(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("Failed to get all orders by user id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.OrderResponse]{
		Code:   200,
		Status: "Success to get all orders by user id",
		Data:   response,
	})
}
