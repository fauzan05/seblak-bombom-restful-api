package middleware

import (
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewRole(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := GetCurrentUser(c)
		if auth.Role != helper.ADMIN {
			userUseCase.Log.Warn("admin access only!")
			return fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}
		return c.Next()
	}
}
