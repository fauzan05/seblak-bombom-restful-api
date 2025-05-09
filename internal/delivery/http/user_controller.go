package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Create(ctx , request)
	if err != nil {
		c.Log.Warnf("failed to register an user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.UserResponse]{
		Code:   201,
		Status: "success to register an user",
		Data:   response,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Authenticate(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to login : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.UserTokenResponse]{
		Code:   200,
		Status: "success to login",
		Data:   response,
	})
}

func (c *UserController) GetCurrent(ctx *fiber.Ctx) error {	
	response := middleware.GetCurrentUser(ctx)

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.UserResponse]{
		Code:   200,
		Status: "success to get user data by token",
		Data:   response,
	})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	// ambil data form update
	dataRequest := new(model.UpdateUserRequest)
	err := ctx.BodyParser(dataRequest)
	if err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}
	// ambil data current user dari auth
	auth := middleware.GetCurrentUser(ctx)

	response, err := c.UseCase.Update(ctx.Context(), dataRequest, auth)
	if err != nil {
		c.Log.Warnf("failed to update user data : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.UserResponse]{
		Code:   200,
		Status: "success to update user data",
		Data:   response,
	})
}

func (c *UserController) UpdatePassword(ctx *fiber.Ctx) error {
	// ambil data form update
	dataRequest := new(model.UpdateUserPasswordRequest)
	err := ctx.BodyParser(dataRequest)
	if err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	// ambil data current user dari auth
	auth := middleware.GetCurrentUser(ctx)
	response, err := c.UseCase.UpdatePassword(ctx.Context(), dataRequest, auth)
	if err != nil {
		c.Log.Warnf("failed to update user password : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to update user password",
		Data:   response,
	})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	tokenRequest := new(model.GetUserByTokenRequest)
	// tangkap token dari header
	result := ctx.GetReqHeaders()
	tokenRequest.Token = result["Authorization"][0]
	response, err := c.UseCase.Logout(ctx.Context(), tokenRequest)
	if err != nil {
		c.Log.Warnf("failed to delete user token : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to logout",
		Data:   response,
	})
}

func (c *UserController) RemoveAccount(ctx *fiber.Ctx) error {
	// ambil data form update
	dataRequest := new(model.DeleteCurrentUserRequest)
	if err := ctx.BodyParser(dataRequest); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	// ambil data current user dari auth
	auth := middleware.GetCurrentUser(ctx)
	response, err := c.UseCase.RemoveCurrentAccount(ctx.Context(), dataRequest, auth)
	if err != nil {
		c.Log.Warnf("failed to delete current user : %+v", err)
		return err
	}
	
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to delete current user",
		Data:   response,
	})
}
