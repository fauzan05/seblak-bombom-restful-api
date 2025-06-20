package middleware

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewAuthUseAdminKey(authConfig *model.AuthConfig, log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		getAuthAppToken := c.Get("X-Admin-Key", "NOT_FOUND")
		fmt.Printf("getAuthAppToken: %s\n & adminCreationKey: %s", getAuthAppToken, authConfig.AdminCreationKey)
		if getAuthAppToken != "NOT_FOUND" && getAuthAppToken != authConfig.AdminCreationKey {
			log.Warnf("admin key isn't valid!")
			return fiber.NewError(fiber.StatusUnauthorized, "admin key isn't valid!")
		}
		return c.Next()
	}
}
