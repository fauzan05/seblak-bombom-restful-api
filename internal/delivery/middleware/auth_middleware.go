package middleware

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &model.GetUserByTokenRequest{
			Token: c.Get("Authorization", "NOT_FOUND"),
		}
		userUseCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUseCase.GetUserByToken(c.Context(), request)
		if err != nil {
			userUseCase.Log.Warnf("Token isn't valid : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUseCase.Log.Debugf("User : %+v", auth.Email)
		c.Locals("auth", auth)
		return c.Next()
	}
}

func GetCurrentUser(ctx *fiber.Ctx) *model.UserResponse {
	return ctx.Locals("auth").(*model.UserResponse)
}
