package config

import (
	"seblak-bombom-restful-api/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName: config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork: config.GetBool("web.prefork"),
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