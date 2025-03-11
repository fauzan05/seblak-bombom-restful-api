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
			userUseCase.Log.Warn("Admin access only!")
			return fiber.NewError(fiber.StatusUnauthorized, "Admin access only!")
		}
		return c.Next()
	}
}
