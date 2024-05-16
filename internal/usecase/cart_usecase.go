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

type CartUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	CartRepository    *repository.CartRepository
	ProductRepository *repository.ProductRepository
}

func NewCartUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	cartRepository *repository.CartRepository, productRepository *repository.ProductRepository) *CartUseCase {
	return &CartUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		CartRepository:    cartRepository,
		ProductRepository: productRepository,
	}
}

func (c *CartUseCase) Add(ctx context.Context, request *model.CreateCartRequest) (*model.CartResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	// dicek apakah produknya ada atau tidak
	newProduct := new(entity.Product)
	newProduct.ID = request.ProductID
	if err := c.ProductRepository.FindById(tx, newProduct); err != nil {
		c.Log.Warnf("Failed to find product by id into product table : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	// cek apakah produk tersedia atau tidak
	if newProduct.Stock < 1 {
		c.Log.Warnf("Product was out of stock : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	// cek apakah permintaan melebihi stok yang tersedia
	newProduct.Stock -= request.Quantity
	if newProduct.Stock < 0 {
		// jika jumlah kuantitasnya melebihi stok yang tersedia
		c.Log.Warnf("Quantity request out of stock from product : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// dicek terlebih dahulu apakah ada cart dengan user yang sama dan produk yang sama.
	newCart := new(entity.Cart)
	cartCount, err := c.CartRepository.FindCartDuplicate(tx, newCart, request.UserID, request.ProductID)
	if err != nil {
		c.Log.Warnf("Failed to find cart duplicate into cart table : %+v", err)
		return nil, fiber.ErrInternalServerError
	}


	if cartCount > 0 {
		// jika ada, maka cukup update saja quantity-nya
		newCart.Quantity += request.Quantity
		newCart.TotalPrice += newProduct.Price * float32(request.Quantity)

		if err := c.CartRepository.Update(tx, newCart); err != nil {
			c.Log.Warnf("Failed to update quantity of product in same the cart : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if request.Quantity < 0 {
			c.Log.Warnf("Quantity must be positive number : %+v", err)
			return nil, fiber.ErrBadRequest
		}

	} else if cartCount < 1 {
		newCart.UserID = request.UserID
		newCart.ProductID = request.ProductID
		newCart.Name = newProduct.Name
		newCart.Quantity = request.Quantity
		newCart.Price = newProduct.Price
		newCart.TotalPrice += newProduct.Price * float32(request.Quantity)

		if err := c.CartRepository.Create(tx, newCart); err != nil {
			c.Log.Warnf("Failed to insert data into cart table : %+v", err)
			return nil, fiber.ErrBadRequest
		}
	}

	// stok produk tersebut kurangi sesuai dengan jumlah quantity request
	if err := c.ProductRepository.Update(tx, newProduct); err != nil {
		c.Log.Warnf("Failed to update stock of product : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CartToResponse(newCart), nil
}
