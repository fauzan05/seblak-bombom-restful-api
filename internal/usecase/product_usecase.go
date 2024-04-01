package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"

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
}

func NewProductUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	categoryRepository *repository.CategoryRepository, productRepository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CategoryRepository: categoryRepository,
		ProductRepository:  productRepository,
	}
}

func (c *ProductUseCase) Add(ctx context.Context, request *model.CreateProductRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
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
		return nil, fiber.ErrInternalServerError
	}

	if err := c.ProductRepository.FindWithJoins(tx, newProduct, "Category"); err != nil {
		c.Log.Warnf("Failed find product by id with preload from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) Get(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	if err := c.ProductRepository.FindWith2Preloads(tx, newProduct, "Category", "Images"); err != nil {
		c.Log.Warnf("Failed get product from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ProductToResponse(newProduct), nil
}

func (c *ProductUseCase) GetAll(ctx context.Context) (*[]model.ProductResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newProducts := new([]entity.Product)
	if err := c.ProductRepository.FindAllWith2Preloads(tx, newProducts, "Category", "Images"); err != nil {
		c.Log.Warnf("Failed get all products from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ProductsToResponse(newProducts), nil
}

func (c *ProductUseCase) Update(ctx context.Context, request *model.UpdateProductRequest) (*model.ProductResponse, error) {
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

	newProduct := new(entity.Product)
	newProduct.ID = request.ID
	if err := c.ProductRepository.Delete(tx, newProduct); err != nil {
		c.Log.Warnf("Failed delete product by id : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
