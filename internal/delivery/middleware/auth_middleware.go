package middleware

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		getToken := c.Cookies("access_token", "NOT_FOUND")
		if getToken == "NOT_FOUND" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		request := &model.GetUserByTokenRequest{
			Token: getToken,
		}

		userUseCase.Log.Debugf("token from cookie: %s", request.Token)
		fmt.Println("token from cookie:", request.Token)
		auth, err := userUseCase.GetUserByToken(c.Context(), request)
		if err != nil {
			userUseCase.Log.Warnf("invalid token: %+v", err)
			return fiber.NewError(fiber.StatusUnauthorized, "token isn't valid")
		}

		c.Locals("auth", auth)
		return c.Next()
	}
}

func GetCurrentUser(ctx *fiber.Ctx) *model.UserResponse {
	return ctx.Locals("auth").(*model.UserResponse)
}
