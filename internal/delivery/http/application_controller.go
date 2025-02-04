package http

import (
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
		c.Log.Warnf("Cannot parse multipart form data: %+v", err)
		return err
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

	request.OpeningHours = form.Value["opening_hours"][0]
	request.ClosingHours = form.Value["closing_hours"][0]
	request.Address = strings.TrimSpace(form.Value["address"][0])
	request.GoogleMapsLink = strings.TrimSpace(form.Value["google_maps_link"][0])
	request.Description = strings.TrimSpace(form.Value["description"][0])
	request.PhoneNumber = strings.TrimSpace(form.Value["phone_number"][0])
	request.Email = strings.TrimSpace(form.Value["email"][0])
	request.InstagramName = strings.TrimSpace(form.Value["instagram_name"][0])
	request.InstagramLink = strings.TrimSpace(form.Value["instagram_link"][0])
	request.TwitterName = strings.TrimSpace(form.Value["twitter_name"][0])
	request.TwitterLink = strings.TrimSpace(form.Value["twitter_link"][0])
	request.FacebookName = strings.TrimSpace(form.Value["facebook_name"][0])
	request.FacebookLink = strings.TrimSpace(form.Value["facebook_link"][0])
	
	response, err := c.UseCase.Add(ctx, request)
	if err != nil {
		c.Log.Warnf("Failed to create new application : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   201,
		Status: "Success to create a new application",
		Data:   response,
	})
}

func (c *ApplicationController) Get(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Get(ctx.Context())
	if err != nil {
		c.Log.Warnf("Failed to get an application : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ApplicationResponse]{
		Code:   200,
		Status: "Success to get an application",
		Data:   response,
	})
}
