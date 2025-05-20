package http

import (
	"encoding/json"
	"fmt"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase/xendit"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type XenditCallbackController struct {
	Log            *logrus.Logger
	UseCase        *usecase.XenditCallbackUseCase
	FrontEndConfig *model.FrontEndConfig
}

func NewXenditCallbackController(useCase *usecase.XenditCallbackUseCase, logger *logrus.Logger, frontEndConfig *model.FrontEndConfig) *XenditCallbackController {
	return &XenditCallbackController{
		Log:            logger,
		UseCase:        useCase,
		FrontEndConfig: frontEndConfig,
	}
}

func (c *XenditCallbackController) GetPaymentRequestCallbacks(ctx *fiber.Ctx) error {
	// Menangkap raw body
	rawBody := ctx.Body()
	var requestData model.XenditGetPaymentRequestCallbackStatus
	err := json.Unmarshal(rawBody, &requestData)
	if err != nil {
		c.Log.Warnf("Failed to unmarshall request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to unmarshall request body : %+v", err))
	}

	lang, _ := requestData.Data.Metadata["lang"].(string)
	if lang == "" {
		lang = "en"
	}
	requestData.Lang = helper.Languange(lang)
	timeZone, _ := requestData.Data.Metadata["time_zone"].(string)
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	requestData.TimeZone = *loc
	requestData.BaseFrontEndURL = c.FrontEndConfig.BaseURL
	err = c.UseCase.UpdateStatusPaymentRequestCallback(ctx, &requestData)
	if err != nil {
		c.Log.Warnf("Failed to process xendit payment request callback : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "Success to get xendit payment request callback",
	})
}

func (c *XenditCallbackController) GetPayoutRequestCallbacks(ctx *fiber.Ctx) error {
	// Menangkap raw body
	rawBody := ctx.Body()
	var requestData model.XenditGetPayoutRequestCallbackStatus
	err := json.Unmarshal(rawBody, &requestData)
	if err != nil {
		c.Log.Warnf("Failed to unmarshall request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to unmarshall request body : %+v", err))
	}

	err = c.UseCase.UpdateStatusPayoutRequestCallback(ctx, &requestData)
	if err != nil {
		c.Log.Warnf("Failed to process xendit payout request callback : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "Success to get xendit payout request callback",
	})
}
