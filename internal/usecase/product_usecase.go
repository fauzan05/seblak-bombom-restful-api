package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"seblak-bombom-restful-api/internal/entity"
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
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if request.Stock < 0 {
		c.Log.Warnf("Stock must be positive number : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newProduct := new(entity.Product)
	newProduct.CategoryId = request.CategoryId
	newProduct.Name = request.Name
	newProduct.Description = request.Description
	newProduct.Price = request.Price
	newProduct.Stock = request.Stock

	if err := c.ProductRepository.Create(tx, newProduct); err != nil {
		c.Log.Warnf("Failed create product into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to store data!")
	}

	if err := c.ProductRepository.FindWithJoins(tx, newProduct, "Category"); err != nil {
		c.Log.Warnf("Failed find product by id with preload from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Couldn't find category because category not created yet! ")
	}

	newImages := make([]entity.Image, len(files))

	for i, file := range files {
		fmt.Printf("File #%d: %s with position %s\n", i+1, file.Filename, positions[i])
		// fmt.Printf("File #%d: %s\n", i+1, file.Filename)
		hashedFilename := hashFileName(file.Filename)
		var position, _ = strconv.Atoi(positions[i])
		// Tambahkan data ke struct ImageAddRequest
		newImages[i].ProductId = newProduct.ID
		newImages[i].FileName = hashedFilename
		newImages[i].Type = file.Header.Get("Content-Type")
		newImages[i].Position = position
		newImages[i].Created_At = time.Now()

		// Simpan file ke direktori uploads
		err := fiberContext.SaveFile(file, fmt.Sprintf("../uploads/images/products/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("Failed to save file: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file!")
		}
	}

	// Simpan gambar ke database
	if err := c.ImageRepository.CreateInBatch(tx, &newImages); err != nil {
		c.Log.Warnf("Failed to save file data into database: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file!")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) Get(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	// Validasi request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID

	// Mengambil data produk
	if err := c.ProductRepository.FindWith2Preloads(c.DB.WithContext(ctx), newProduct, "Category", "Images"); err != nil {
		c.Log.Warnf("Failed get product from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Mengembalikan response produk
	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) GetAll(ctx context.Context, page int, perPage int) (*[]model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page <= 0 {
		page = 1
	}
	
	newProducts := new([]entity.Product)
	if err := c.ProductRepository.FindAllWith2Preloads(tx, newProducts, "Category", "Images", page, perPage); err != nil {
		c.Log.Warnf("Failed get all products from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ProductsToResponse(newProducts), nil
}

func (c *ProductUseCase) Update(ctx context.Context, fiberContext *fiber.Ctx, request *model.UpdateProductRequest, newImageFiles []*multipart.FileHeader, newImagePositions []string, updateCurrentImages model.UpdateImagesRequest, deletedImages model.DeleteImagesRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	newProduct.CategoryId = request.CategoryId
	newProduct.Name = request.Name
	newProduct.Description = request.Description
	newProduct.Price = request.Price
	newProduct.Stock = request.Stock

	if err := c.ProductRepository.Update(tx, newProduct); err != nil {
		c.Log.Warnf("Failed update product by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Jika user menambahkan gambar baru
	if len(newImageFiles) > 0 {
		newImages := make([]entity.Image, len(newImageFiles))

		for i, file := range newImageFiles {
			fmt.Printf("File #%d: %s with position %s\n", i+1, file.Filename, newImagePositions[i])
			// fmt.Printf("File #%d: %s\n", i+1, file.Filename)
			hashedFilename := hashFileName(file.Filename)
			var position, _ = strconv.Atoi(newImagePositions[i])
			// Tambahkan data ke struct ImageAddRequest
			newImages[i].ProductId = newProduct.ID
			newImages[i].FileName = hashedFilename
			newImages[i].Type = file.Header.Get("Content-Type")
			newImages[i].Position = position
			newImages[i].Created_At = time.Now()

			// Simpan file ke direktori uploads
			err := fiberContext.SaveFile(file, fmt.Sprintf("../uploads/images/products/%s", hashedFilename))
			if err != nil {
				c.Log.Warnf("Failed to save file: %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file!")
			}
		}

		// Simpan gambar baru ke database
		if err := c.ImageRepository.CreateInBatch(tx, &newImages); err != nil {
			c.Log.Warnf("Failed to save file data into database: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file!")
		}
	}

	if len(updateCurrentImages.Images) > 0 {
		for _, updateCurrentImage := range updateCurrentImages.Images {
			updateImage := entity.Image{
				ID: updateCurrentImage.ID,
			}

			if err := c.ImageRepository.FindById(tx, &updateImage); err != nil {
				c.Log.Warnf("Failed to find images by id : %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			updateImage.Position = updateCurrentImage.Position
			// Perbarui posisi gambar
			if err := c.ImageRepository.Update(tx, &updateImage); err != nil {
				c.Log.Warnf("Failed to save updated current image data into database: %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update images!")
			}
		}
	}

	if len(deletedImages.Images) > 0 {
		for _, deletedImage := range deletedImages.Images {
			deleteImage := entity.Image{
				ID: deletedImage.ID,
			}

			if err := c.ImageRepository.Delete(tx, &deleteImage); err != nil {
				c.Log.Warnf("Failed to delete current image data into database: %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete images!")
			}
		}
	}

	if err := c.ProductRepository.FindWithJoins(tx, newProduct, "Category"); err != nil {
		c.Log.Warnf("Failed get product from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) Delete(ctx context.Context, request *model.DeleteProductRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	currentImage := new([]entity.Image)
	if err := c.ImageRepository.FindImageByProductId(tx, currentImage, request.ID); err != nil {
		c.Log.Warnf("Failed find product images by product id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to find product images!")
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	if err := c.ProductRepository.Delete(tx, newProduct); err != nil {
		c.Log.Warnf("Failed delete product by id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete products!")
	}

	newImage := new(entity.Image)
	if err := c.ImageRepository.DeleteByProductId(tx, newImage, request.ID); err != nil {
		c.Log.Warnf("Failed to delete image file data from database: %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete products image!")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

// Fungsi untuk membuat hash dari nama file
func hashFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d-%s", timestamp, originalName)))
	return fmt.Sprintf("%x", hash[:8]) + filepath.Ext(originalName) // Menggunakan 8 karakter pertama hash
}
