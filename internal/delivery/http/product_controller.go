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

type ProductController struct {
	Log     *logrus.Logger
	UseCase *usecase.ProductUseCase
}

func NewProductController(useCase *usecase.ProductUseCase, logger *logrus.Logger) *ProductController {
	return &ProductController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	// Buat direktori uploads jika belum ada
	if _, err := os.Stat("../uploads/images/products/"); os.IsNotExist(err) {
		os.MkdirAll("../uploads/images/products/", os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.Warnf("Cannot parse multipart form data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse multipart form data : %+v", err))
	}

	request := new(model.CreateProductRequest)
	categoryID, _ := strconv.ParseUint(form.Value["category_id"][0], 10, 64)
	request.CategoryId = categoryID
	request.Name = form.Value["name"][0]
	request.Description = form.Value["description"][0]
	parsePrice64, _ := strconv.ParseFloat(form.Value["price"][0], 64)
	request.Price = float32(parsePrice64)
	request.Stock, _ = strconv.Atoi(form.Value["stock"][0])

	files := form.File["images"]
	positions := form.Value["positions"]

	response, err := c.UseCase.Add(ctx.Context(), ctx, request, files, positions)
	if err != nil {
		c.Log.Warnf("Failed to create a new product : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   201,
		Status: "Success to create a new product",
		Data:   response,
	})
}

func (c *ProductController) Get(ctx *fiber.Ctx) error {
	getId := ctx.Params("productId")
	productId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("Failed to convert product_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to convert product_id to integer : %+v", err))
	}

	productRequest := new(model.GetProductRequest)
	productRequest.ID = uint64(productId)
	response, err := c.UseCase.Get(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to get selected product : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   200,
		Status: "Success to get selected product",
		Data:   response,
	})
}

func (c *ProductController) GetAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search", "")
	trimSearch := strings.TrimSpace(search)

	getCategoryId := ctx.Query("category_id", "");
	var categoryId uint64
	categoryId = 0
	if getCategoryId != "" {
		getValueConvert, err := strconv.ParseUint(ctx.Query("category_id", ""), 10, 64)
		if err != nil {
			c.Log.Warnf("Invalid or missing category_id : %+v", err)
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid or missing category_id : %+v", err))
		}
		categoryId = getValueConvert
	}

	// ambil data sorting
	getColumn := ctx.Query("column", "");
	getSortBy := ctx.Query("sort_by", "desc");

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

	response, totalProducts, totalPages, err := c.UseCase.GetAll(ctx.Context(), page, perPage, trimSearch, categoryId, getColumn, getSortBy)
	if err != nil {
		c.Log.Warnf("Failed to find all products : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.ProductResponse]{
		Code:   200,
		Status: "Success to get all products",
		Data:   response,
		TotalDatas: totalProducts,
		TotalPages: totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *ProductController) Edit(ctx *fiber.Ctx) error {
	// Buat direktori uploads jika belum ada
	if _, err := os.Stat("../uploads/images/products/"); os.IsNotExist(err) {
		os.MkdirAll("../uploads/images/products/", os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.Warnf("Cannot parse multipart form data: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Cannot parse multipart form data: %+v", err))
	}

	request := new(model.UpdateProductRequest)
	getProductId, err := strconv.ParseUint(ctx.Params("productId"), 10, 64)
	if err != nil {
		c.Log.Warnf("Invalid product ID : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid product ID : %+v", err))
	}
	request.ID = getProductId
	categoryID, err := strconv.ParseUint(form.Value["category_id"][0], 10, 64)
	if err != nil {
		c.Log.Warnf("Invalid category ID : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid category ID : %+v", err))
	}
	request.CategoryId = categoryID
	request.Name = form.Value["name"][0]
	request.Description = form.Value["description"][0]
	parsePrice64, err := strconv.ParseFloat(form.Value["price"][0], 64)
	if err != nil {
		c.Log.Warnf("Invalid price : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid price : %+v", err))
	}
	request.Price = float32(parsePrice64)
	request.Stock, err = strconv.Atoi(form.Value["stock"][0])
	if err != nil {
		c.Log.Warnf("Invalid stock : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid stock : %+v", err))
	}

	// Inisialisasi NEW IMAGES
	newImageFiles := form.File["new_images"]
	newImagePositions := form.Value["new_positions"]

	// Inisialisasi CURRENT IMAGES
	updateImagesRequest := model.UpdateImagesRequest{}
	for i, imageId := range form.Value["current_images"] {
		imageId, err := strconv.ParseUint(imageId, 10, 64)
		if err != nil {
			c.Log.Warnf("Invalid image ID : %+v", err)
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid image ID : %+v", err))
		}
		currentPosition, err := strconv.Atoi(form.Value["current_positions"][i])
		if err != nil {
			c.Log.Warnf("Invalid current position : %+v", err)
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid current position : %+v", err))
		}
		currentImage := model.ImageUpdateRequest{
			ID:       imageId,
			Position: currentPosition,
		}
		updateImagesRequest.Images = append(updateImagesRequest.Images, currentImage)
	}

	// Inisialisasi DELETED IMAGES
	deleteImagesRequest := model.DeleteImagesRequest{}
	if len(form.Value["images_deleted"]) > 0 {
		imagesDeleted := strings.Split(form.Value["images_deleted"][0], ",")
		for _, imageId := range imagesDeleted {
			imageId, err := strconv.ParseUint(imageId, 10, 64)
			if err != nil {
				c.Log.Warnf("Invalid image ID : %+v", err)
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid image ID : %+v", err))
			}
			deleteImage := model.DeleteImageRequest{
				ID: imageId,
			}
			deleteImagesRequest.Images = append(deleteImagesRequest.Images, deleteImage)
		}
	}

	response, err := c.UseCase.Update(ctx.Context(), ctx, request, newImageFiles, newImagePositions, updateImagesRequest, deleteImagesRequest)
	if err != nil {
		c.Log.Warnf("Failed to update an product by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   200,
		Status: "Success to update an product by id",
		Data:   response,
	})
}

func (c *ProductController) Remove(ctx *fiber.Ctx) error {
	idsParam := ctx.Query("ids")
	if idsParam == "" {
		c.Log.Warnf("Parameter 'ids' is required")
		return fiber.NewError(fiber.StatusBadRequest, "Parameter 'ids' is required")
	}

	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var productIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.Log.Warnf("Parameter 'ids' is required")
				return fiber.NewError(fiber.StatusBadRequest, "Parameter 'ids' is required")
			}
			productIds = append(productIds, id)
		}
	}

	productRequest := new(model.DeleteProductRequest)
	productRequest.IDs = productIds

	response, err := c.UseCase.Delete(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to delete selected products : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete selected products",
		Data:   response,
	})
}
