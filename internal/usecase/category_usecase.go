package usecase

import (
	"context"
	"fmt"
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

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	categoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CategoryRepository: categoryRepository,
	}
}

func (c *CategoryUseCase) Add(ctx context.Context, request *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCategory := new(entity.Category)
	newCategory.Name = request.Name
	newCategory.Description = request.Description
	if err := c.CategoryRepository.Create(tx, newCategory); err != nil {
		c.Log.Warnf("Failed create category into database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetById(ctx context.Context, request *model.GetCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request query params : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	if err := c.CategoryRepository.FindById(tx, newCategory); err != nil {
		c.Log.Warnf("Can't find category by id : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.CategoryResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page <= 0 {
		page = 1
	}

	var result []map[string]interface{} // entity kosong yang akan diisi
	if err := c.CategoryRepository.FindCategoriesPagination(tx, &result, page, perPage, search, sortingColumn, sortBy,); err != nil {
		c.Log.Warnf("Failed to find all categories : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	newCategories := new([]entity.Category)
	err := MapCategories(result, newCategories)
	if err != nil {
		c.Log.Warnf("Failed map categories : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	var totalPages int = 0
	getAllCategories := new(entity.Category)
	totalCategories, err := c.CategoryRepository.CountCategoryItems(tx, getAllCategories, search)
	if err != nil {
		c.Log.Warnf("Failed to count categories: %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	// Hitung total halaman
	totalPages = int(totalCategories / int64(perPage))
	if totalCategories%int64(perPage) > 0 {
		totalPages++
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	return converter.CategoriesToResponse(newCategories), totalCategories, totalPages, nil
}

func (c *CategoryUseCase) Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	newCategory.Name = request.Name
	newCategory.Description = request.Description
	if err := c.CategoryRepository.Update(tx, newCategory); err != nil {
		c.Log.Warnf("Failed to update category : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) Delete(ctx context.Context, request *model.DeleteCategoryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newCategories := []entity.Category{}

	for _, idProduct := range request.IDs {
		newCategory := entity.Category{
			ID: idProduct,
		}

		newCategories = append(newCategories, newCategory)
	}

	if err := c.CategoryRepository.DeleteInBatch(tx, &newCategories); err != nil {
		c.Log.Warnf("Failed to delete category : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func MapCategories(rows []map[string]interface{}, results *[]entity.Category) (error) {
	layoutWithZone := "2006-01-02T15:04:05-07:00"

	for _, row := range rows {
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
		categoryCreatedAt, err := time.Parse(layoutWithZone, categoryCreatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse category_created_at: %v", err)
		}

		categoryUpdatedAtStr, _ := row["category_updated_at"].(string)
		categoryUpdatedAt, err := time.Parse(layoutWithZone, categoryUpdatedAtStr)
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

		// Tambahkan ke hasil
		*results = append(*results, newCategory)
	}

	return nil
}
