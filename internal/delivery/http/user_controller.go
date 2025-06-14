package http

import (
	"fmt"
	"html/template"
	"os"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log            *logrus.Logger
	UseCase        *usecase.UserUseCase
	AuthConfig     *model.AuthConfig
	FrontEndConfig *model.FrontEndConfig
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger, authConfig *model.AuthConfig, frontEndConfig *model.FrontEndConfig) *UserController {
	return &UserController{
		Log:            logger,
		UseCase:        useCase,
		AuthConfig:     authConfig,
		FrontEndConfig: frontEndConfig,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	request.TimeZone = *loc
	if request.Role == enum_state.ADMIN {
		adminKey := ctx.Get("X-Admin-Key", "")
		if adminKey != c.AuthConfig.AdminCreationKey {
			c.Log.Warnf("invalid admin creation key!")
			return fiber.NewError(fiber.StatusForbidden, "invalid admin creation key!")
		}
	}

	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	response, err := c.UseCase.Create(ctx, request)
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

func (c *UserController) VerifyEmailRegistration(ctx *fiber.Ctx) error {
	getVerifyToken := ctx.Params("token", "")
	request := new(model.VerifyEmailRegisterRequest)
	request.VerificationToken = getVerifyToken
	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.TimeZone = *loc
	request.BaseFrontEndURL = c.FrontEndConfig.BaseURL
	response, err := c.UseCase.VerifyEmailRegistration(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to verify email registration an user : %+v", err)
		return err
	}

	url := fmt.Sprintf("/verified-success/%s?lang=%s&email=%s", getVerifyToken, getLang, response.Email)
	return ctx.Redirect(url, fiber.StatusFound)
}

func (c *UserController) ShowVerifiedSuccess(ctx *fiber.Ctx) error {
	ctx.Type("html", "utf-8")
	bodyBuilder := new(strings.Builder)
	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	getVerifyToken := ctx.Params("token", "")
	getEmail := ctx.Query("email", "")
	err := c.UseCase.ValidateVerifyTokenIsValid(ctx, getVerifyToken, getEmail)
	if err != nil {
		c.Log.Warnf("%+v", err)
		htmlBytes, _ := os.ReadFile(fmt.Sprintf("internal/templates/%s/pages/internal_error.html", getLang))
		return ctx.Status(500).Type("html").SendString(string(htmlBytes))
	}

	templatePath := fmt.Sprintf("../internal/templates/%s/pages/verified_success.html", getLang)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	loginUrl := fmt.Sprintf("%s/login", c.FrontEndConfig.BaseURL)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"LoginURL": loginUrl,
		"Email":    getEmail,
	})

	if err != nil {
		c.Log.Warnf("failed to execute verified_success page : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return ctx.SendString(bodyBuilder.String())
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
	request := new(model.DeleteCurrentUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.TimeZone = *loc
	// ambil data current user dari auth
	auth := middleware.GetCurrentUser(ctx)
	response, err := c.UseCase.RemoveCurrentAccount(ctx.Context(), request, auth)
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

func (c *UserController) CreateForgotPassword(ctx *fiber.Ctx) error {
	request := new(model.CreateForgotPassword)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	response, err := c.UseCase.AddForgotPassword(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to create an forgot password request : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.PasswordResetResponse]{
		Code:   200,
		Status: "success to create an forgot password request",
		Data:   response,
	})
}

func (c *UserController) ValidateForgotPassword(ctx *fiber.Ctx) error {
	getId := ctx.Params("passwordResetId")
	passwordResetId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert password_reset_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert password_reset_id to integer : %+v", err))
	}

	request := new(model.ValidateForgotPassword)
	request.ID = uint64(passwordResetId)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.ValidateForgotPassword(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to validate an forgot password request : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to validate an forgot password request",
		Data:   response,
	})
}

func (c *UserController) ResetPassword(ctx *fiber.Ctx) error {
	getId := ctx.Params("passwordResetId")
	passwordResetId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert password_reset_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert password_reset_id to integer : %+v", err))
	}

	request := new(model.PasswordResetRequest)
	request.ID = uint64(passwordResetId)
	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.TimeZone = *loc
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Reset(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to validate an forgot password request : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to validate an forgot password request",
		Data:   response,
	})
}
