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
			userUseCase.Log.Warn("Just an admin can access")
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}
}
