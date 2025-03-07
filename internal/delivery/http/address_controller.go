package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Log     *logrus.Logger
	UseCase *usecase.AddressUseCase
}

func NewAddressController(useCase *usecase.AddressUseCase, logger *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     logger,
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
		Code:   201,
		Status: "Success to add new address",
		Data:   response,
	})
}

func (c *AddressController) GetAll(ctx *fiber.Ctx) error {
	auth := middleware.GetCurrentUser(ctx)
	response, err := c.UseCase.GetAll(auth)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.AddressResponse]{
		Code:   200,
		Status: "Success to get all address by current user",
		Data:   response,
	})
}

func (c *AddressController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("addressId")
	addressId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert address id : %+v", err)
		return err
	}
	addressRequest := &model.GetAddressRequest{
		ID: uint64(addressId),
	}

	response, err := c.UseCase.GetById(ctx.Context(), addressRequest)
	if err != nil {
		c.Log.Warnf("Failed to find address by id : %+v", err)
		return err
	}
	
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.AddressResponse]{
		Code:   200,
		Status: "Success to get an address",
		Data:   response,
	})
}

func (c *AddressController) Update(ctx *fiber.Ctx) error {
	getId := ctx.Params("addressId")
	addressId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert address id : %+v", err)
		return err
	}
	// ambil data dari body
	addressRequest := new(model.UpdateAddressRequest)
	err = ctx.BodyParser(addressRequest)
	if err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	// ambil data current user dari auth
	auth := middleware.GetCurrentUser(ctx)

	addressRequest.ID = uint64(addressId)
	addressRequest.UserId = auth.ID

	response, err := c.UseCase.Edit(ctx.Context(), addressRequest)
	if err != nil {
		c.Log.Warnf("Failed to edit address : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.AddressResponse]{
		Code:   200,
		Status: "Success to update an address",
		Data:   response,
	})
}

func (c *AddressController) Remove(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"error": "Parameter 'ids' is required",
		})
	}
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var addressIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if idStr != "" {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid ID: %s", idStr),
				})
			}
			addressIds = append(addressIds, id)
		}
	}

	deleteAddress := new(model.DeleteAddressRequest)
	deleteAddress.IDs = addressIds

	response, err := c.UseCase.Delete(ctx.Context(), deleteAddress)
	if err != nil {
		c.Log.Warnf("Failed to delete address : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete an address",
		Data:   response,
	})
}
