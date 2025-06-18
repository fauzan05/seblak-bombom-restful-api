package config

import (
	"seblak-bombom-restful-api/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName: config.GetString("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
		Prefork: config.GetBool("WEB_PREFORK"),
		BodyLimit: 100 * 1024 * 1024,
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		helper.SaveToLogError(err.Error())
		return c.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}	
}