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
		c.Log.Warnf("Cannot parse multipart form data: %+v", err)
		return err
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
		c.Log.Warnf("Failed to create new product : %+v", err)
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
		c.Log.Warnf("Failed to convert product id : %+v", err)
		return err
	}
	productRequest := new(model.GetProductRequest)
	productRequest.ID = uint64(productId)

	response, err := c.UseCase.Get(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("Failed to find product by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.ProductResponse]{
		Code:   200,
		Status: "Success to get an product",
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
			// Handle kasus ketika category_id tidak valid atau kosong
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or missing category_id",
			})
		}
		categoryId = getValueConvert
	}

	// ambil data sorting
	getColumn := ctx.Query("column", "");
	getSortBy := ctx.Query("sort_by", "desc");


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
		return err
	}

	request := new(model.UpdateProductRequest)
	getProductId, _ := strconv.ParseUint(ctx.Params("productId"), 10, 64)
	request.ID = getProductId
	categoryID, _ := strconv.ParseUint(form.Value["category_id"][0], 10, 64)
	request.CategoryId = categoryID
	request.Name = form.Value["name"][0]
	request.Description = form.Value["description"][0]
	parsePrice64, _ := strconv.ParseFloat(form.Value["price"][0], 64)
	request.Price = float32(parsePrice64)
	request.Stock, _ = strconv.Atoi(form.Value["stock"][0])

	// Inisialisasi NEW IMAGES
	newImageFiles := form.File["new_images"]
	newImagePositions := form.Value["new_positions"]

	// Inisialisasi CURRENT IMAGES
	updateImagesRequest := model.UpdateImagesRequest{}
	for i, imageId := range form.Value["current_images"] {
		imageId, _ := strconv.ParseUint(imageId, 10, 64)
		currentPosition, _ := strconv.Atoi(form.Value["current_positions"][i])
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
			imageId, _ := strconv.ParseUint(imageId, 10, 64)
			deleteImage := model.DeleteImageRequest{
				ID: imageId,
			}
			deleteImagesRequest.Images = append(deleteImagesRequest.Images, deleteImage)
		}
	}

	response, err := c.UseCase.Update(ctx.Context(), ctx, request, newImageFiles, newImagePositions, updateImagesRequest, deleteImagesRequest)
	if err != nil {
		c.Log.Warnf("Failed to update product by id : %+v", err)
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
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"error": "Parameter 'ids' is required",
		})
	}

	// Pisahkan string menjadi array menggunakan koma sebagai delimiter
	idStrings := strings.Split(idsParam, ",")
	var productIds []uint64

	// Konversi setiap elemen menjadi integer
	for _, idStr := range idStrings {
		if (idStr != "") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Invalid ID: %s", idStr),
				})
			}
			productIds = append(productIds, id)
		}
	}

	productRequest := new(model.DeleteProductRequest)
	productRequest.IDs = productIds

	response, err := c.UseCase.Delete(ctx.Context(), productRequest)
	if err != nil {
		c.Log.Warnf("%+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[bool]{
		Code:   200,
		Status: "Success to delete products by ids",
		Data:   response,
	})
}
