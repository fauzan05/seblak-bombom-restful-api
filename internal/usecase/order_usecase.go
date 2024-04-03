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

type OrderUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	OrderRepository    *repository.OrderRepository
	ProductRepository  *repository.ProductRepository
	CategoryRepository *repository.CategoryRepository
	AddressRepository  *repository.AddressRepository
	DiscountRepository *repository.DiscountRepository
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountRepository) *OrderUseCase {
		return &OrderUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		OrderRepository:    orderRepository,
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
		AddressRepository:  addressRepository,
		DiscountRepository: discountRepository,
	}
}

func (c *OrderUseCase) Add(ctx context.Context, request *model.CreateOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newOrder := new(entity.Order)
	newProduct := new(entity.Product)
	// temukan produk untuk memastikan ketersediaan
	newProduct.ID = request.ProductId
	count, err := c.ProductRepository.FindAndCountById(tx, newProduct)
	if err != nil {
		c.Log.Warnf("Failed to get product by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count < 1 {
		c.Log.Warnf("Find product by id not found : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if newProduct.Stock < 1 {
		c.Log.Warnf("Product out of stock : %+v", err)
		return nil, fiber.ErrNotFound
	}
	
	newOrder.ProductId = newProduct.ID
	newOrder.ProductName = newProduct.Name
	newOrder.ProductDescription = newProduct.Description
	newOrder.Price = newProduct.Price
	newOrder.Quantity = request.Quantity
	newOrder.Amount = newProduct.Price * newOrder.Quantity

	if request.DiscountCode != "" {
		newDiscount := new(entity.Discount)
		count, err := c.DiscountRepository.CountDiscountByCode(tx, newDiscount, request.DiscountCode) 
		if err != nil {
			c.Log.Warnf("Can't find discount by code : %+v", err)
			return nil, fiber.ErrNotFound
		}

		// cek apakah diskonnya ada dan statusnya aktif (true)
		if count > 0 && newDiscount.Status {
			// cek apakah sudah kadaluarsa atau belum
			if newDiscount.End.Before(time.Now()) {
				if newDiscount.Type == helper.PERCENT {
					newOrder.Amount -= newDiscount.Value / 100
				} else {
					newOrder.Amount -= newDiscount.Value
				}
			}
		}
	}

	// user/customer data
	newOrder.UserId = request.UserId
	newOrder.FirstName = request.FirstName
	newOrder.LastName = request.LastName
	newOrder.Email = request.Email
	newOrder.Phone = request.Phone

	// payment
	newOrder.PaymentMethod = request.PaymentMethod
	newOrder.PaymentStatus = request.PaymentStatus
	// newOrder.

	newOrder.ProductName = request.ProductName
	newOrder.ProductDescription = request.ProductDescription
	newOrder.Price = request.Price
	newOrder.Quantity = request.Quantity
	newOrder.Amount = request.Amount
	newOrder.DiscountValue = request.Amount
	newOrder.DiscountType = request.DiscountType
	newOrder.UserId = request.UserId
	newOrder.FirstName = request.FirstName
	newOrder.LastName = request.LastName
	newOrder.Email = request.Email
	newOrder.Phone = request.Phone
	newOrder.PaymentMethod = request.PaymentMethod
	newOrder.PaymentStatus = request.PaymentStatus
	newOrder.DeliveryStatus = request.DeliveryStatus
	newOrder.IsDelivery = request.IsDelivery
	newOrder.DeliveryCost = request.DeliveryCost
	newOrder.CategoryName = request.CategoryName
	newOrder.CompleteAddress = request.CompleteAddress
	newOrder.GoogleMapLink = request.GoogleMapLink
	newOrder.Distance = request.Distance

	if err := c.OrderRepository.Create(tx, newOrder); err != nil {
		c.Log.Warnf("failed to create new order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	invoice := fmt.Sprintf("INV/%d/USER/%d", newOrder.ID, newOrder.UserId)
	newOrder.Invoice = invoice
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("failed to add invoice code : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrderToResponse(newOrder), nil
}
