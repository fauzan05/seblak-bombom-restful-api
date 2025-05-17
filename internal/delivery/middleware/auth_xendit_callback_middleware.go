package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewAuthXenditCallback(config *viper.Viper, log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		xenditCallbackToken := config.GetString("xendit.test.callback_token")
		requestToken := c.Get("X-Callback-Token", "NOT_FOUND")
		if requestToken != "NOT_FOUND" && requestToken != xenditCallbackToken {
			log.Warnf("xendit callback token isn't valid!")
			return fiber.NewError(fiber.StatusUnauthorized, "xendit callback token isn't valid!")
		}
		c.Locals("xendit_callback_token", xenditCallbackToken)
		return c.Next()
	}
}
