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
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newCategory := new(entity.Category)
	newCategory.Name = request.Name
	newCategory.Description = request.Description
	if err := c.CategoryRepository.Create(tx, newCategory); err != nil {
		c.Log.Warnf("Failed to create category into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed create category into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetById(ctx context.Context, request *model.GetCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	if err := c.CategoryRepository.FindById(tx, newCategory); err != nil {
		c.Log.Warnf("Can't find category by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Can't find category by id : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.CategoryResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	var result []map[string]any // entity kosong yang akan diisi
	if err := c.CategoryRepository.FindCategoriesPagination(tx, &result, page, perPage, search, sortingColumn, sortBy); err != nil {
		c.Log.Warnf("Failed to find all categories : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find all categories : %+v", err))
	}

	newCategories := new([]entity.Category)
	err := MapCategories(result, newCategories)
	if err != nil {
		c.Log.Warnf("Failed to map categories : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed map categories : %+v", err))
	}

	var totalPages int = 0
	newCategory := new(entity.Category)
	totalCategories, err := c.CategoryRepository.CountCategoryItems(tx, newCategory, search)
	if err != nil {
		c.Log.Warnf("Failed to count categories : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to count categories: %+v", err))
	}

	// Hitung total halaman
	totalPages = int(totalCategories / int64(perPage))
	if totalCategories%int64(perPage) > 0 {
		totalPages++
	}

	return converter.CategoriesToResponse(newCategories), totalCategories, totalPages, nil
}

func (c *CategoryUseCase) Update(ctx context.Context, request *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	newCategory.Name = request.Name
	newCategory.Description = request.Description
	newCategory.Updated_At = time.Now().UTC()
	if err := c.CategoryRepository.Update(tx, newCategory); err != nil {
		c.Log.Warnf("Failed to update category : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update category : %+v", err))
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) Delete(ctx context.Context, request *model.DeleteCategoryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
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
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to delete category : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return true, nil
}

func MapCategories(rows []map[string]any, results *[]entity.Category) error {

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

		// Tambahkan ke hasil
		*results = append(*results, newCategory)
	}

	return nil
}
