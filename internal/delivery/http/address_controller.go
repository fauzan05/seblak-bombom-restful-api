package http

import (
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Log *logrus.Logger
	UseCase *usecase.AddressUseCase
}

func NewAddressController(useCase *usecase.AddressUseCase, logger *logrus.Logger) *AddressController {
	return &AddressController{
		Log: logger,
		UseCase: useCase,
	}
}

func (c *AddressController) Add(ctx *fiber.Ctx) error {
	request := new(model.AddressCreateRequest)
	token := new(model.GetUserByTokenRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	getToken := ctx.GetReqHeaders()
	token.Token = getToken["Authorization"][0]

	response, err := c.UseCase.Create(ctx.Context(), request, token)
	if err != nil {
		c.Log.Warnf("Failed to create address : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.AddressResponse]{
		Code: 201,
		Status: "Success to add new address",
		Data: response,
	})

}