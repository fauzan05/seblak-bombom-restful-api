package http

import (
	"fmt"
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
		c.Log.Warnf("Failed to convert order id : %+v", err)
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
		Status: "Success get xendit transaction by order id",
		Data:   response,
	})
}

func (c *XenditQRCodeTransctionController) GetCallbacks(ctx *fiber.Ctx) error {
	// Menangkap semua header
	rawHeader := ctx.GetReqHeaders()
	for key, values := range rawHeader {
		fmt.Printf("Header: %s = %s\n", key, values)
	}
	getCallbackToken := rawHeader["X-Callback-Token"]
	fmt.Println("TOKEN : ", getCallbackToken)

	// Menangkap raw body
	rawBody := ctx.Body()
	fmt.Println("Raw Body:", string(rawBody))

	// Struct untuk menampung data JSON dari request body
	var requestData map[string]interface{}

	// Parsing body JSON ke dalam struct
	if err := ctx.BodyParser(&requestData); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "Success to get xendit transaction callback",
		"data":   requestData,
	})
}
