package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MidtransCoreAPIOrderController struct {
	Log     *logrus.Logger
	UseCase *usecase.MidtransCoreAPIOrderUseCase
}

func NewMidtransCoreAPIOrderController(useCase *usecase.MidtransCoreAPIOrderUseCase, logger *logrus.Logger) *MidtransCoreAPIOrderController {
	return &MidtransCoreAPIOrderController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *MidtransCoreAPIOrderController) CreateCoreAPI(ctx *fiber.Ctx) error {
	snapRequest := new(model.CreateMidtransCoreAPIOrderRequest)
	if err := ctx.BodyParser(snapRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), snapRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new midtrans core api order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.MidtransCoreAPIOrderResponse]{
		Code:   201,
		Status: "Success to create a new midtrans core api order",
		Data:   response,
	})
}

func (c *MidtransCoreAPIOrderController) GetCoreAPIOrderNotification(ctx *fiber.Ctx) error {
	coreApiRequest := new(model.GetMidtransCoreAPIOrderRequest)
	getId := ctx.Params("orderId")
	orderId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order id : %+v", err)
		return err
	}
	
	coreApiRequest.OrderId = uint64(orderId)
	response, err := c.UseCase.Get(ctx.Context(), coreApiRequest)
	if err != nil {
		c.Log.Warnf("Failed to get a notification midtrans order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.MidtransCoreAPIOrderResponse]{
		Code:   200,
		Status: "Success to get a midtrans notification order",
		Data:   response,
	})
}
