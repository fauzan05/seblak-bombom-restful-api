package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime/multipart"
	"os"
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

func (c *ProductUseCase) GetAll(ctx context.Context, page int, perPage int, search string, categoryId uint64, sortingColumn string, sortBy string) (*[]model.ProductResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page <= 0 {
		page = 1
	}

	var result []map[string]interface{} // entity kosong yang akan diisi
	if err := c.ProductRepository.GetProductsWithPagination(tx, &result, page, perPage, search, sortingColumn, sortBy, categoryId); err != nil {
		c.Log.Warnf("Failed get all products from database : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}
	
	newProducts := new([]entity.Product)
	err := MapProducts(result, newProducts)
	if err != nil {
		c.Log.Warnf("Failed map products : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}
	
	var totalPages int = 0
	getAllProducts := new(entity.Product)
	totalProducts, err := c.ProductRepository.CountProductItems(tx, getAllProducts, search, categoryId)
	if err != nil {
		c.Log.Warnf("Failed to count products: %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	// Hitung total halaman
	totalPages = int(totalProducts / int64(perPage))
	if totalProducts%int64(perPage) > 0 {
		totalPages++
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	return converter.ProductsToResponse(newProducts), totalProducts, totalPages, nil
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

	currentImages := new([]entity.Image)
	if err := c.ImageRepository.FindImagesByProductIds(tx, currentImages, request.IDs); err != nil {
		c.Log.Warnf("Failed find product images by product id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to find product images!")
	}

	// Mengecek apakah gambar kosong atau nil
	if len(*currentImages) != 0 {
		// hapus semua gambar di database
		if err := c.ImageRepository.DeleteInBatch(tx, currentImages); err != nil {
			c.Log.Warnf("Failed to delete product images in database : %+v", err)
			return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete product images in database!")
		}

		filePath := "../uploads/images/products/"

		// Hapus file gambar
		for _, currentImage := range *currentImages {
			err = os.Remove(filePath + currentImage.FileName)
			if err != nil {
				fmt.Printf("Error deleting file: %v\n", err)
				return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete product images in directory!")
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
		c.Log.Warnf("Failed delete in batch product by id : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete in batch products!")
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

func MapProducts(rows []map[string]interface{}, results *[]entity.Product) (error) {

	for _, row := range rows {
		// fmt.Printf("Produk Id : %s | DATANYA : %d\n", row["product_id"], len(row["images"].([]map[string]interface{})))
		imagesConvertToInterface, ok := row["images"].([]map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse images field, expected []map[string]interface{} but got %T", row["images"])
		}

		newImages := make([]entity.Image, 0)
		for _, image := range imagesConvertToInterface {
			// Ambil dan validasi image_id
			imageIdStr, ok := image["image_id"].(string)
			if !ok || imageIdStr == "" {
				return fmt.Errorf("missing or invalid image_id")
			}
			imageId, err := strconv.ParseUint(imageIdStr, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse image_id : %v", err)
			}

			imageFilename, _ := image["image_filename"].(string)
			imagePosition, _ := strconv.Atoi(image["image_position"].(string))
			imageType, _ := image["image_type"].(string)

			imageCreatedAtStr, _ := image["image_created_at"].(string)
			imageCreatedAt, err := time.Parse(time.RFC3339, imageCreatedAtStr)
			if err != nil {
				return fmt.Errorf("failed to parse image_created_at: %v", err)
			}

			imageUpdatedAtStr, _ := image["image_created_at"].(string)
			imageUpdatedAt, err := time.Parse(time.RFC3339, imageUpdatedAtStr)
			if err != nil {
				return fmt.Errorf("failed to parse image_created_at: %v", err)
			}

			newImage := entity.Image{
				ID:         imageId,
				FileName:   imageFilename,
				Position:   imagePosition,
				Type:       imageType,
				Created_At: imageCreatedAt,
				Updated_At: imageUpdatedAt,
			}

			newImages = append(newImages, newImage)
		}

		// Ambil dan validasi category_id
		categoryIdStr, ok := row["category_id"].(string)
		if !ok || categoryIdStr == "" {
			return fmt.Errorf("missing or invalid category_id")
		}

		categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse category_id: %v", err)
		}

		// Ambil field kategori
		categoryName, _ := row["category_name"].(string)
		categoryDesc, _ := row["category_desc"].(string)

		// Parse created_at dan updated_at kategori
		categoryCreatedAtStr, _ := row["category_created_at"].(string)
		categoryCreatedAt, err := time.Parse(time.RFC3339, categoryCreatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse category_created_at: %v", err)
		}

		categoryUpdatedAtStr, _ := row["category_updated_at"].(string)
		categoryUpdatedAt, err := time.Parse(time.RFC3339, categoryUpdatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse category_updated_at: %v", err)
		}

		// Buat objek kategori
		newCategory := entity.Category{
			ID:          categoryId,
			Name:        categoryName,
			Description: categoryDesc,
			Created_At:  categoryCreatedAt,
			Updated_At:  categoryUpdatedAt,
		}

		// Ambil dan validasi product_id
		productIdStr, ok := row["product_id"].(string)
		if !ok || productIdStr == "" {
			return fmt.Errorf("missing or invalid product_id")
		}
		productId, err := strconv.ParseUint(productIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse product_id: %v", err)
		}

		// Ambil field produk
		productName, _ := row["product_name"].(string)
		productDesc, _ := row["product_desc"].(string)

		// Parse harga dan stok produk
		productPriceStr, _ := row["product_price"].(string)
		productPrice, err := strconv.ParseFloat(productPriceStr, 32)
		if err != nil {
			return fmt.Errorf("failed to parse product_price: %v", err)
		}
		productStockStr, _ := row["product_stock"].(string)
		productStock, err := strconv.Atoi(productStockStr)
		if err != nil {
			return fmt.Errorf("failed to parse product_stock: %v", err)
		}

		// Parse created_at dan updated_at produk
		productCreatedAtStr, _ := row["product_created_at"].(string)
		productCreatedAt, err := time.Parse(time.RFC3339, productCreatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse product_created_at: %v", err)
		}
		productUpdatedAtStr, _ := row["product_updated_at"].(string)
		productUpdatedAt, err := time.Parse(time.RFC3339, productUpdatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse product_updated_at: %v", err)
		}

		// Buat objek produk
		newProduct := entity.Product{
			ID:          productId,
			Name:        productName,
			Description: productDesc,
			Price:       float32(productPrice),
			Stock:       productStock,
			Created_At:  productCreatedAt,
			Updated_At:  productUpdatedAt,
			Category:    &newCategory,
			Images:      newImages,
		}

		// Tambahkan ke hasil
		*results = append(*results, newProduct)
	}

	return nil
}
