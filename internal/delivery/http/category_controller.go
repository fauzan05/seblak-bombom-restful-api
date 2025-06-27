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
	if _, err := os.Stat("uploads/images/categories/"); os.IsNotExist(err) {
		os.MkdirAll("uploads/images/categories/", os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.Warnf("cannot parse multipart form data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse multipart form data : %+v", err))
	}

	request := new(model.CreateCategoryRequest)
	request.Name = strings.TrimSpace(form.Value["name"][0])
	request.Description = strings.TrimSpace(form.Value["description"][0])
	if len(form.File["image"]) > 0 {
		request.Image = form.File["image"][0]
	} else {
		request.Image = nil
	}

	response, err := c.UseCase.Add(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to create new category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   201,
		Status: "success to create a new category",
		Data:   response,
	})
}

func (c *CategoryController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("categoryId")
	categoryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert category_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}
	categoryRequest := new(model.GetCategoryRequest)
	categoryRequest.ID = uint64(categoryId)

	response, err := c.UseCase.GetById(ctx.Context(), categoryRequest)
	if err != nil {
		c.Log.Warnf("failed to find category by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   200,
		Status: "success to get an category",
		Data:   response,
	})
}

func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search", "")
	trimSearch := strings.TrimSpace(search)

	// ambil data sorting
	getColumn := ctx.Query("column", "")
	getSortBy := ctx.Query("sort_by", "desc")
	isActive := ctx.Query("is_active", "")

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

	response, totalProducts, totalPages, err := c.UseCase.GetAll(ctx.Context(), page, perPage, trimSearch, getColumn, getSortBy, isActive)
	if err != nil {
		c.Log.Warnf("failed to find all categories : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.CategoryResponse]{
		Code:   200,
		Status: "success to get all category",
		Data:   response,
		TotalDatas: totalProducts,
		TotalPages: totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *CategoryController) Edit(ctx *fiber.Ctx) error {
	if _, err := os.Stat("uploads/images/categories/"); os.IsNotExist(err) {
		os.MkdirAll("uploads/images/categories/", os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.Warnf("cannot parse multipart form data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse multipart form data : %+v", err))
	}

	getId := ctx.Params("categoryId")
	categoryId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert category_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert category_id to integer : %+v", err))
	}

	request := new(model.UpdateCategoryRequest)
	request.ID = uint64(categoryId)
	request.Name = strings.TrimSpace(form.Value["name"][0])
	request.Description = strings.TrimSpace(form.Value["description"][0])
	if len(form.File["image"]) > 0 {
		request.Image = form.File["image"][0]
	} else {
		request.Image = nil
	}

	response, err := c.UseCase.Update(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to update category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.CategoryResponse]{
		Code:   200,
		Status: "success to update category",
		Data:   response,
	})
}

func (c *CategoryController) Remove(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		c.Log.Warnf("parameter 'ids' is required")
		return fiber.NewError(fiber.StatusBadRequest, "parameter 'ids' is required")
	}
	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var categoryIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.Log.Warnf("invalid category ID : %s", err)
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid category ID : %s", err))
			}
			categoryIds = append(categoryIds, id)
		}
	}

	deleteCategory := new(model.DeleteCategoryRequest)
	deleteCategory.IDs = categoryIds

	response, err := c.UseCase.Delete(ctx.Context(), deleteCategory)
	if err != nil {
		c.Log.Warnf("failed to delete category : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "success to delete category",
		Data:   response,
	})
}
