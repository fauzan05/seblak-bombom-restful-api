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

type CategoryController struct {
	Log     *logrus.Logger
	UseCase *usecase.CategoryUseCase
}

func NewCategoryController(useCase *usecase.CategoryUseCase, logger *logrus.Logger) *CategoryController {
	return &CategoryController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateCategoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}

	response, err := c.UseCase.Add(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create new category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   201,
		Status: "Success to create a new category",
		Data:   response,
	})
}

func (c *CategoryController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("categoryId")
	categoryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert category id : %+v", err)
		return err
	}
	categoryRequest := new(model.GetCategoryRequest)
	categoryRequest.ID = uint64(categoryId)

	response, err := c.UseCase.GetById(ctx.Context(), categoryRequest)
	if err != nil {
		c.Log.Warnf("Failed to find category by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   200,
		Status: "Success to get an category",
		Data:   response,
	})
}

func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
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

	response, totalProducts, totalPages, err := c.UseCase.GetAll(ctx.Context(), page, perPage, trimSearch, getColumn, getSortBy)
	if err != nil {
		c.Log.Warnf("Failed to find all categories : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.CategoryResponse]{
		Code:   200,
		Status: "Success to get all category",
		Data:   response,
		TotalDatas: totalProducts,
		TotalPages: totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *CategoryController) Edit(ctx *fiber.Ctx) error {
	getId := ctx.Params("categoryId")
	categoryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert category id : %+v", err)
		return err
	}

	updateCategory := new(model.UpdateCategoryRequest)
	if err := ctx.BodyParser(updateCategory); err != nil {
		c.Log.Warnf("Cannot parse data : %+v", err)
		return err
	}
	updateCategory.ID = uint64(categoryId)
	response, err := c.UseCase.Update(ctx.Context(), updateCategory)
	if err != nil {
		c.Log.Warnf("Failed to edit category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   200,
		Status: "Success to update category",
		Data:   response,
	})
}

func (c *CategoryController) Remove(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"error": "Parameter 'ids' is required",
		})
	}
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var categoryIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid ID: %s", idStr),
				})
			}
			categoryIds = append(categoryIds, id)
		}
	}

	deleteCategory := new(model.DeleteCategoryRequest)
	deleteCategory.IDs = categoryIds

	response, err := c.UseCase.Delete(ctx.Context(), deleteCategory)
	if err != nil {
		c.Log.Warnf("Failed to delete category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete category",
		Data:   response,
	})
}
