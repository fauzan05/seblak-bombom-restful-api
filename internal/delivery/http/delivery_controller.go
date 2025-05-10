package http

import (
	"fmt"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DeliveryController struct {
	Log     *logrus.Logger
	UseCase *usecase.DeliveryUseCase
}

func NewDeliveryController(useCase *usecase.DeliveryUseCase, logger *logrus.Logger) *DeliveryController {
	return &DeliveryController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *DeliveryController) Create(ctx  *fiber.Ctx) error {
	request := new(model.CreateDeliveryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("failed to create/update delivery setting : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.DeliveryResponse]{
		Code:   201,
		Status: "success to create a delivery settings",
		Data:   response,
	})
}

func (c *DeliveryController) GetAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search", "")
	trimSearch := strings.TrimSpace(search)

	// ambil data sorting
	getColumn := ctx.Query("column", "")
	getSortBy := ctx.Query("sort_by", "desc")

	// Ambil query parameter 'per_page' dengan default value 10 jika tidak disediakan
	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil {
		c.Log.Warnf("invalid 'per_page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid 'per_page' parameter : %+v", err))
	}

	// Ambil query parameter 'page' dengan default value 1 jika tidak disediakan
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		c.Log.Warnf("invalid 'page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid 'page' parameter : %+v", err))
	}

	response, totalDeliveries, totalPages, err := c.UseCase.GetAll(ctx.Context(), page, perPage, trimSearch, getColumn, getSortBy)
	if err != nil {
		c.Log.Warnf("failed to get a delivery setting data : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.DeliveryResponse]{
		Code:   200,
		Status: "success to get a delivery settings",
		Data:   response,
		TotalDatas: totalDeliveries,
		TotalPages: totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *DeliveryController) Update(ctx  *fiber.Ctx) error {
	deliveryRequest := new(model.UpdateDeliveryRequest)
	if err := ctx.BodyParser(deliveryRequest); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}

	getId := ctx.Params("deliveryId")
	deliveryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert delivery_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert delivery_id to integer : %+v", err))
	}

	deliveryRequest.ID = uint64(deliveryId)
	response, err := c.UseCase.Edit(ctx.Context(), deliveryRequest)
	if err != nil {
		c.Log.Warnf("failed to update delivery setting : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DeliveryResponse]{
		Code:   200,
		Status: "success to update selected delivery settings",
		Data:   response,
	})
}

func (c *DeliveryController) Remove(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		c.Log.Warnf("parameter 'ids' is required")
		return fiber.NewError(fiber.StatusBadRequest, "parameter 'ids' is required")
	}
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var deliveryIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.Log.Warnf("invalid delivery ID : %s", err)
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid delivery ID : %s", err))
			}
			deliveryIds = append(deliveryIds, id)
		}
	}

	deleteDelivery := new(model.DeleteDeliveryRequest)
	deleteDelivery.IDs = deliveryIds
	response, err := c.UseCase.Delete(ctx.Context(), deleteDelivery)
	if err != nil {
		c.Log.Warnf("failed to remove selected delivery settings : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to remove selected delivery settings",
		Data:   response,
	})
}