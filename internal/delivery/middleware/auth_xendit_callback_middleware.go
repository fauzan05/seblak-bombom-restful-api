package middleware

import (

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


func NewAuthXenditCallback(config viper.Viper, log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		xenditCallbackToken := config.GetString("xendit.test.callback_token")
		requestToken := c.Get("X-Callback-Token", "NOT_FOUND")
		log.Debugf("X-Callback-Token : %s", requestToken)

		if requestToken != xenditCallbackToken {
			log.Warnf("Xendit callback token isn't valid!")
			return fiber.ErrUnauthorized
		}
		c.Locals("xendit_callback_token", xenditCallbackToken)
		return c.Next()
	}
}