package http

import (
	"fmt"
	"os"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ApplicationController struct {
	Log     *logrus.Logger
	UseCase *usecase.ApplicationUseCase
}

func NewApplicationController(useCase *usecase.ApplicationUseCase, logger *logrus.Logger) *ApplicationController {
	return &ApplicationController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ApplicationController) Create(ctx *fiber.Ctx) error {
	// Buat direktori uploads jika belum ada
	if _, err := os.Stat("../uploads/images/application/"); os.IsNotExist(err) {
		os.MkdirAll("../uploads/images/application/", os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.Warnf("cannot parse multipart form data: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse multipart form data: %+v", err))
	}

	request := new(model.CreateApplicationRequest)
	if len(form.Value["id"]) > 0 {
		request.ID, _ = strconv.ParseUint(form.Value["id"][0], 10, 64)
	}

	request.AppName = strings.TrimSpace(form.Value["app_name"][0])
	if len(form.File["logo_filename"]) > 0 {
		request.Logo = form.File["logo_filename"][0]
	} else {
		request.Logo = nil
	}

	getFirst := func(key string) string {
		values, ok := form.Value[key]
		if !ok || len(values) == 0 {
			return "" // atau kamu bisa return default value, atau error
		}
		return strings.TrimSpace(values[0])
	}

	request.OpeningHours = getFirst("opening_hours")
	request.ClosingHours = getFirst("closing_hours")
	request.Address = getFirst("address")
	request.GoogleMapsLink = getFirst("google_maps_link")
	request.Description = getFirst("description")
	request.PhoneNumber = getFirst("phone_number")
	request.Email = getFirst("email")
	request.InstagramName = getFirst("instagram_name")
	request.InstagramLink = getFirst("instagram_link")
	request.TwitterName = getFirst("twitter_name")
	request.TwitterLink = getFirst("twitter_link")
	request.FacebookName = getFirst("facebook_name")
	request.FacebookLink = getFirst("facebook_link")

	response, err := c.UseCase.Add(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to create new application : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   200,
		Status: "success to create/update application settings",
		Data:   response,
	})
}

func (c *ApplicationController) Get(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Get(ctx.Context())
	if err != nil {
		c.Log.Warnf("failed to get application settings : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   200,
		Status: "success to get application settings",
		Data:   response,
	})
}
