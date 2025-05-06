package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CategoryRepository *repository.CategoryRepository
	ProductRepository  *repository.ProductRepository
	ImageRepository    *repository.ImageRepository
}

func NewProductUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	categoryRepository *repository.CategoryRepository, productRepository *repository.ProductRepository, imageRepository *repository.ImageRepository) *ProductUseCase {
	return &ProductUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CategoryRepository: categoryRepository,
		ProductRepository:  productRepository,
		ImageRepository:    imageRepository,
	}
}

func (c *ProductUseCase) Add(ctx context.Context, fiberContext *fiber.Ctx, request *model.CreateProductRequest, files []*multipart.FileHeader, positions []string) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	if request.Stock < 0 {
		c.Log.Warnf("stock must be positive number!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "stock must be positive number!")
	}

	if request.Price < 0 {
		c.Log.Warnf("price must be positive number!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "price must be positive number!")
	}

	if len(positions) == 0 {
		c.Log.Warnf("image position must be included!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "image position must be included!")
	}

	if len(files) == 0 {
		c.Log.Warnf("images must be uploaded!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "images must be uploaded!")
	}

	if len(files) != len(positions) {
		c.Log.Warnf("each uploaded image must have a corresponding position!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "each uploaded image must have a corresponding position!")
	}

	if len(files) > 5 {
		c.Log.Warnf("you can upload up to 5 images only!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "you can upload up to 5 images only!")
	}

	// cek apakah catgory-nya ada
	newCategory := new(entity.Category)
	newCategory.ID = request.CategoryId
	count, err := c.CategoryRepository.FindAndCountById(tx, newCategory)
	if err != nil {
		c.Log.Warnf("failed to find category from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find category from database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("category is not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "category is not found!")
	}

	newProduct := new(entity.Product)
	newProduct.CategoryId = request.CategoryId
	newProduct.Name = request.Name
	newProduct.Description = request.Description
	newProduct.Price = request.Price
	newProduct.Stock = request.Stock

	if err := c.ProductRepository.Create(tx, newProduct); err != nil {
		c.Log.Warnf("failed to create product into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create product into database : %+v", err))
	}

	if err := c.ProductRepository.FindWithJoins(tx, newProduct, "Category"); err != nil {
		c.Log.Warnf("failed to find product by id from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find product by id from database : %+v", err))
	}

	newImages := make([]entity.Image, len(files))

	for i, file := range files {
		err = helper.ValidateFile(1, file)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// fmt.Printf("File #%d: %s\n", i+1, file.Filename)
		hashedFilename := hashFileName(file.Filename)
		var position, _ = strconv.Atoi(positions[i])
		// Tambahkan data ke struct ImageAddRequest
		newImages[i].ProductId = newProduct.ID
		newImages[i].FileName = hashedFilename
		newImages[i].Type = file.Header.Get("Content-Type")
		newImages[i].Position = position

		// Simpan file ke direktori uploads
		err := fiberContext.SaveFile(file, fmt.Sprintf("../uploads/images/products/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("failed to save uploaded file : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save uploaded file : %+v", err))
		}
	}

	// Simpan gambar ke database
	if err := c.ImageRepository.CreateInBatch(tx, &newImages); err != nil {
		c.Log.Warnf("failed to save images file_name into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save images file_name into database : %+v", err))
	}

	if err := c.ProductRepository.FindWith2Preloads(tx, newProduct, "Category", "Images"); err != nil {
		c.Log.Warnf("failed to get product from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get product from database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}
	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) Get(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	// Validasi request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	count, err := c.ProductRepository.FindAndCountById(tx, newProduct)
	if err != nil {
		c.Log.Warnf("failed to find product from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find product from database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("product is not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "product is not found!")
	}

	// Mengambil data produk
	if err := c.ProductRepository.FindWith2Preloads(tx, newProduct, "Category", "Images"); err != nil {
		c.Log.Warnf("failed to get product by id from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get product by id from database : %+v", err))
	}

	// Mengembalikan response produk
	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) GetAll(ctx context.Context, page int, perPage int, search string, categoryId uint64, sortingColumn string, sortBy string) (*[]model.ProductResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	if sortingColumn == "" {
		sortingColumn = "products.id"
	}

	newPagination := new(repository.Pagination)
	newPagination.Page = page
	newPagination.PageSize = perPage
	newPagination.Column = sortingColumn
	newPagination.SortBy = sortBy
	allowedColumns := map[string]bool{
		"products.id":          true,
		"products.name":        true,
		"categories.name":      true,
		"products.description": true,
		"products.created_at":  true,
		"products.updated_at":  true,
	}

	if !allowedColumns[newPagination.Column] {
		c.Log.Warnf("invalid sort column : %s", newPagination.Column)
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid sort column : %s", newPagination.Column))
	}
	
	products, totalProduct, err := repository.Paginate(tx, &entity.Product{}, newPagination, func(d *gorm.DB) *gorm.DB {
		return d.Joins("JOIN categories ON categories.id = products.category_id").
			Preload("Category").
			Preload("Images").Where("products.name LIKE ?", "%"+search+"%")
	})

	if err != nil {
		c.Log.Warnf("failed to paginate category : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to paginate category : %+v", err))
	}

	// Hitung total halaman
	var totalPages int = 0
	totalPages = int(totalProduct / int64(perPage))
	if totalProduct%int64(perPage) > 0 {
		totalPages++
	}

	return converter.ProductsToResponse(&products), totalProduct, totalPages, nil
}

func (c *ProductUseCase) Update(ctx context.Context, fiberContext *fiber.Ctx, request *model.UpdateProductRequest, newImageFiles []*multipart.FileHeader, newImagePositions []string, updateCurrentImages model.UpdateImagesRequest, deletedImages model.DeleteImagesRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	if request.Stock < 0 {
		c.Log.Warnf("stock must be positive number!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "stock must be positive number!")
	}

	if request.Price < 0 {
		c.Log.Warnf("price must be positive number!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "price must be positive number!")
	}

	// cek apakah jumlah new image file itu sama dengan jumlah new image position
	if len(newImageFiles) > 0 && len(newImageFiles) != len(newImagePositions) {
		c.Log.Warnf("each new uploaded image must have a corresponding position!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "each new uploaded image must have a corresponding position!")
	}

	// cek apakah gambarnya melebihi 5
	if (len(newImageFiles) + len(updateCurrentImages.Images)) > 5 {
		c.Log.Warnf("you can upload up to 5 images only!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "you can upload up to 5 images only!")
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.CategoryId
	count, err := c.CategoryRepository.FindAndCountById(tx, newCategory)
	if err != nil {
		c.Log.Warnf("failed to find category from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find category from database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("category is not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "category is not found!")
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	count, err = c.ProductRepository.FindAndCountById(tx, newProduct)
	if err != nil {
		c.Log.Warnf("failed to find product from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find product from database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("product is not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "product is not found!")
	}

	newProduct.CategoryId = request.CategoryId
	newProduct.Name = request.Name
	newProduct.Description = request.Description
	newProduct.Price = request.Price
	newProduct.Stock = request.Stock

	if err := c.ProductRepository.Update(tx, newProduct); err != nil {
		c.Log.Warnf("failed update product by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed update product by id : %+v", err))
	}

	// Jika user menambahkan gambar baru
	if len(newImageFiles) > 0 {
		newImages := make([]entity.Image, len(newImageFiles))

		for i, file := range newImageFiles {
			err = helper.ValidateFile(1, file)
			if err != nil {
				c.Log.Warnf(err.Error())
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}

			// fmt.Printf("File #%d: %s\n", i+1, file.Filename)
			hashedFilename := hashFileName(file.Filename)
			var position, _ = strconv.Atoi(newImagePositions[i])
			// Tambahkan data ke struct ImageAddRequest
			newImages[i].ProductId = newProduct.ID
			newImages[i].FileName = hashedFilename
			newImages[i].Type = file.Header.Get("Content-Type")
			newImages[i].Position = position
			newImages[i].CreatedAt = time.Now()

			// Simpan file ke direktori uploads
			err := fiberContext.SaveFile(file, fmt.Sprintf("../uploads/images/products/%s", hashedFilename))
			if err != nil {
				c.Log.Warnf("failed to save uploaded file : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save uploaded file : %+v", err))
			}
		}

		// Simpan gambar baru ke database
		if err := c.ImageRepository.CreateInBatch(tx, &newImages); err != nil {
			c.Log.Warnf("failed to save images file_name into database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save images file_name into database : %+v", err))
		}
	}

	if len(updateCurrentImages.Images) > 0 {
		for _, updateCurrentImage := range updateCurrentImages.Images {
			updateImage := entity.Image{
				ID: updateCurrentImage.ID,
			}

			if err := c.ImageRepository.FindById(tx, &updateImage); err != nil {
				c.Log.Warnf("failed to find images by id : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find images by id : %+v", err))
			}

			updateImage.Position = updateCurrentImage.Position
			// Perbarui posisi gambar
			if err := c.ImageRepository.Update(tx, &updateImage); err != nil {
				c.Log.Warnf("failed to save updated current image data into database : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save updated current image data into database : %+v", err))
			}
		}
	}

	if len(deletedImages.Images) > 0 {
		filePath := "../uploads/images/products/"

		for _, deletedImage := range deletedImages.Images {
			deleteImage := entity.Image{
				ID: deletedImage.ID,
			}

			if err := c.ImageRepository.FindById(tx, &deleteImage); err != nil {
				c.Log.Warnf("failed to find current image data in the database : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find current image data in the database : %+v", err))
			}

			if err := c.ImageRepository.Delete(tx, &deleteImage); err != nil {
				c.Log.Warnf("failed to delete current image data in the database : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete current image data in the database : %+v", err))
			}

			err = os.Remove(filePath + deleteImage.FileName)
			if err != nil {
				fmt.Printf("failed to delete image file : %v\n", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete image file : %v\n", err))
			}
		}
	}

	if err := c.ProductRepository.FindWith2Preloads(tx, newProduct, "Category", "Images"); err != nil {
		c.Log.Warnf("failed to get product from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get product from database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) Delete(ctx context.Context, request *model.DeleteProductRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	currentImages := new([]entity.Image)
	if err := c.ImageRepository.FindImagesByProductIds(tx, currentImages, request.IDs); err != nil {
		c.Log.Warnf("failed to find product images by product id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find product images by product id : %+v", err))
	}

	// Mengecek apakah gambar kosong atau nil
	if len(*currentImages) != 0 {
		// hapus semua gambar di database
		if err := c.ImageRepository.DeleteInBatch(tx, currentImages); err != nil {
			c.Log.Warnf("failed to delete product images in the database : %+v", err)
			return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete product images in the database : %+v", err))
		}

		filePath := "../uploads/images/products/"

		// Hapus file gambar
		for _, currentImage := range *currentImages {
			if currentImage.FileName != "" {
				err = os.Remove(filePath + currentImage.FileName)
				if err != nil {
					fmt.Printf("failed to delete image file : %v\n", err)
					return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete image file : %v\n", err))
				}
			}
		}
	}

	// hapus semua produk
	newProducts := []entity.Product{}

	for _, idProduct := range request.IDs {
		newProduct := entity.Product{
			ID: idProduct,
		}

		newProducts = append(newProducts, newProduct)
	}

	// hapus produk di database
	if err := c.ProductRepository.DeleteInBatch(tx, &newProducts); err != nil {
		c.Log.Warnf("failed delete in batch product by id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed delete in batch product by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}

// Fungsi untuk membuat hash dari nama file
func hashFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256(fmt.Appendf(nil, "%d-%s", timestamp, originalName))
	return fmt.Sprintf("%x", hash[:8]) + filepath.Ext(originalName) // Menggunakan 8 karakter pertama hash
}
