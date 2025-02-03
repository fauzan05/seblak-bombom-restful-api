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
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to create add new discount : %+v", err)
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
		c.Log.Warnf("Invalid 'per_page' parameter")
		return err
	}

	// Ambil query parameter 'page' dengan default value 1 jika tidak disediakan
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		c.Log.Warnf("Invalid 'page' parameter")
		return err
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
		c.Log.Warnf("Failed to convert discount id : %+v", err)
		return err
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
		c.Log.Warnf("Failed to convert discount id : %+v", err)
		return err
	}
	discountRequest := new(model.UpdateDiscountCouponRequest)
	discountRequest.ID = uint64(discountId)
	if err := ctx.BodyParser(discountRequest); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	response, err := c.UseCase.Edit(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to edit discount by id : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.DiscountCouponResponse]{
		Code:   200,
		Status: "Success to edit discount by id",
		Data:   response,
	})
}

func(c *DiscountCouponController) Delete(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"error": "Parameter 'ids' is required",
		})
	}
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var discountCouponIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid ID: %s", idStr),
				})
			}
			discountCouponIds = append(discountCouponIds, id)
		}
	}

	discountRequest := new(model.DeleteDiscountCouponRequest)
	discountRequest.IDs = discountCouponIds

	response, err := c.UseCase.Remove(ctx.Context(), discountRequest)
	if err != nil {
		c.Log.Warnf("Failed to delete discount by id : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete discount by id",
		Data:   response,
	})
}
