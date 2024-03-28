package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log: logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to register new user : %+v", err)
		return err
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.UserResponse]{
		Code: 201,
		Status: "Success to register an user",
		Data: response,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequst)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to login : %+v", err)
		return err
	}

	response, err := c.UseCase.Login(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to login : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.UserTokenResponse]{
		Code: 200,
		Status: "Success to login",
		Data: response,
	})
}

func (c *UserController) GetCurrent(ctx *fiber.Ctx) error {
	request := new(model.GetUserByTokenRequest)
	// tangkap token dari header
	result := ctx.GetReqHeaders()
	request.Token = result["Authorization"][0]

	response, err := c.UseCase.GetUserByToken(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.UserResponse]{
		Code: 200,
		Status: "Success to get user data",
		Data: response,
	})
}