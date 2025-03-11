package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase/xendit"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type XenditQRCodeTransctionController struct {
	Log                 *logrus.Logger
	UseCase             *usecase.XenditTransactionQRCodeUseCase
	
}

func NewXenditQRCodeTransctionController(useCase *usecase.XenditTransactionQRCodeUseCase, logger *logrus.Logger) *XenditQRCodeTransctionController {
	return &XenditQRCodeTransctionController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *XenditQRCodeTransctionController) Create(ctx *fiber.Ctx) error {
	xenditRequest := new(model.CreateXenditTransaction)
	if err := ctx.BodyParser(xenditRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx, xenditRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new xendit transaction order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.XenditTransactionResponse]{
		Code:   201,
		Status: "Success to create a new xendit transaction",
		Data:   response,
	})
}

func (c *XenditQRCodeTransctionController) GetTransaction(ctx *fiber.Ctx) error {
	getId := ctx.Params("orderId")
	orderId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert order_id into integer : %+v", err)
		return err
	}

	xenditRequest := new(model.GetXenditQRCodeTransaction)
	xenditRequest.OrderId = uint64(orderId)

	response, err := c.UseCase.GetTransaction(ctx, xenditRequest)
	if err != nil {
		c.Log.Warnf("Failed to get xendit QR code transaction : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.XenditTransactionResponse]{
		Code:   200,
		Status: "Success to get xendit transaction by order id",
		Data:   response,
	})
}