package http

import (
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

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
		return err
	}
	auth := middleware.GetCurrentUser(ctx)
	orderRequest.UserId = auth.ID
	orderRequest.FirstName = auth.FirstName
	orderRequest.LastName = auth.LastName
	orderRequest.Email = auth.Email
	orderRequest.Phone = auth.Phone
	// filter address yang dimana address is_main = true
	for _, address := range auth.Addresses {
		if address.IsMain {
			orderRequest.CompleteAddress = address.CompleteAddress
			orderRequest.GoogleMapLink = address.GoogleMapLink
			break
		}
	}
	
	response, err := c.UseCase.Add(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   201,
		Status: "Success to create a new order",
		Data:   response,
	})
}