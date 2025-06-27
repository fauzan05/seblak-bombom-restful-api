package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
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

func (c *CategoryUseCase) Add(ctx *fiber.Ctx, request *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	hashedFilename := ""
	if request.ImageFilename != nil {
		err := helper.ValidateFile(1, request.ImageFilename)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		hashedFilename = helper.HashFileName(request.ImageFilename.Filename)
		err = ctx.SaveFile(request.ImageFilename, fmt.Sprintf("uploads/images/categories/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("failed to save category image file : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save category image file : %+v", err))
		}
	}

	newCategory := new(entity.Category)
	newCategory.Name = request.Name
	newCategory.Description = request.Description
	if request.ImageFilename != nil {
		newCategory.ImageFilename = hashedFilename
	}

	if err := c.CategoryRepository.Create(tx, newCategory); err != nil {
		c.Log.Warnf("failed to create category into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed create category into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetById(ctx context.Context, request *model.GetCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	if err := c.CategoryRepository.FindById(tx, newCategory); err != nil {
		c.Log.Warnf("can't find category by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't find category by id : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.CategoryResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	if sortingColumn == "" {
		sortingColumn = "categories.id"
	}

	newPagination := new(repository.Pagination)
	newPagination.Page = page
	newPagination.PageSize = perPage
	newPagination.Column = sortingColumn
	newPagination.SortBy = sortBy
	allowedColumns := map[string]bool{
		"categories.id":          true,
		"categories.name":        true,
		"categories.description": true,
		"categories.created_at":  true,
		"categories.updated_at":  true,
	}

	if !allowedColumns[newPagination.Column] {
		c.Log.Warnf("invalid sort column : %s", newPagination.Column)
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid sort column : %s", newPagination.Column))
	}
	
	categories, totalCategory, err := repository.Paginate(tx, &entity.Category{}, newPagination, func(d *gorm.DB) *gorm.DB {
		return d.Where("categories.name LIKE ?", "%"+search+"%")
	})

	if err != nil {
		c.Log.Warnf("failed to paginate categories : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to paginate categories : %+v", err))
	}

	// Hitung total halaman
	var totalPages int = 0
	totalPages = int(totalCategory / int64(perPage))
	if totalCategory%int64(perPage) > 0 {
		totalPages++
	}

	return converter.CategoriesToResponse(&categories), totalCategory, totalPages, nil
}

func (c *CategoryUseCase) Update(ctx *fiber.Ctx, request *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newCategory := new(entity.Category)
	newCategory.ID = request.ID
	// temukan category apakah ada
	count, err := c.CategoryRepository.FindAndCountById(tx, newCategory)
	if err != nil {
		c.Log.Warnf("failed to find category by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find category by id : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("category not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "category not found!")
	}

	hashedFilename := ""
	if request.ImageFilename != nil {
		err := helper.ValidateFile(1, request.ImageFilename)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		hashedFilename = helper.HashFileName(request.ImageFilename.Filename)
		err = ctx.SaveFile(request.ImageFilename, fmt.Sprintf("uploads/images/categories/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("failed to save category image file : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save category image file : %+v", err))
		}
	}

	newCategory.Name = request.Name
	newCategory.Description = request.Description
	if request.ImageFilename != nil {
		newCategory.ImageFilename = hashedFilename
	}
	newCategory.UpdatedAt = time.Now().UTC()
	if err := c.CategoryRepository.Update(tx, newCategory); err != nil {
		c.Log.Warnf("failed to update category : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update category : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.CategoryToResponse(newCategory), nil
}

func (c *CategoryUseCase) Delete(ctx context.Context, request *model.DeleteCategoryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newCategories := []entity.Category{}
	for _, idProduct := range request.IDs {
		newCategory := entity.Category{
			ID: idProduct,
		}

		newCategories = append(newCategories, newCategory)
	}

	if err := c.CategoryRepository.DeleteInBatch(tx, &newCategories); err != nil {
		c.Log.Warnf("failed to delete category : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete category : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}
