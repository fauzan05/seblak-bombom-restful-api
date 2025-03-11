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

type DiscountCouponController struct {
	Log     *logrus.Logger
	UseCase *usecase.DiscountCouponUseCase
}

func NewDiscountCouponController(useCase *usecase.DiscountCouponUseCase, logger *logrus.Logger) *DiscountCouponController {
	return &DiscountCouponController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *DiscountCouponController) Create(ctx *fiber.Ctx) error {
	discountRequest := new(model.CreateDiscountCouponRequest)
	if err := ctx.BodyParser(discountRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Add(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to create new discount : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.DiscountCouponResponse]{
		Code:   201,
		Status: "Success to create a new discount",
		Data:   response,
	})
}

func (c *DiscountCouponController) GetAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search", "")
	trimSearch := strings.TrimSpace(search)

	// ambil data sorting
	getColumn := ctx.Query("column", "")
	getSortBy := ctx.Query("sort_by", "desc")

	// Ambil query parameter 'per_page' dengan default value 10 jika tidak disediakan
	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil {
		c.Log.Warnf("Invalid 'per_page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid 'per_page' parameter : %+v", err))
	}

	// Ambil query parameter 'page' dengan default value 1 jika tidak disediakan
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		c.Log.Warnf("Invalid 'page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid 'page' parameter : %+v", err))
	}

	response, totalDiscountCoupons, totalPages, err := c.UseCase.GetAll(ctx.Context(), page, perPage, trimSearch, getColumn, getSortBy)
	if err != nil {
		c.Log.Warnf("Failed to get all discounts : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.DiscountCouponResponse]{
		Code:   200,
		Status: "Success to get all discounts",
		Data:   response,
		TotalDatas: totalDiscountCoupons,
		TotalPages: totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *DiscountCouponController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("discountId")
	discountId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert discount_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert discount_id to integer : %+v", err))
	}

	getDiscount := new(model.GetDiscountCouponRequest)
	getDiscount.ID = uint64(discountId)
	response, err := c.UseCase.GetById(ctx.Context(), getDiscount)
	if err != nil {
		c.Log.Warnf("Failed to get discount by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DiscountCouponResponse]{
		Code:   200,
		Status: "Success to get discount by id",
		Data:   response,
	})
}

func (c *DiscountCouponController) Update(ctx *fiber.Ctx) error {
	getId := ctx.Params("discountId")
	discountId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert discount_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert discount_id to integer : %+v", err))
	}

	discountRequest := new(model.UpdateDiscountCouponRequest)
	discountRequest.ID = uint64(discountId)
	if err := ctx.BodyParser(discountRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse data : %+v", err))
	}

	response, err := c.UseCase.Edit(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to update selected discount : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DiscountCouponResponse]{
		Code:   200,
		Status: "Success to update selected discount",
		Data:   response,
	})
}

func(c *DiscountCouponController) Delete(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		c.Log.Warnf("Parameter 'ids' is required")
		return fiber.NewError(fiber.StatusBadRequest, "Parameter 'ids' is required")
	}
	
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var discountCouponIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.Log.Warnf("Invalid ID : %s", err)
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid ID : %s", err))
			}
			discountCouponIds = append(discountCouponIds, id)
		}
	}

	discountRequest := new(model.DeleteDiscountCouponRequest)
	discountRequest.IDs = discountCouponIds
	response, err := c.UseCase.Remove(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to delete selected discount : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete selected discount",
		Data:   response,
	})
}
