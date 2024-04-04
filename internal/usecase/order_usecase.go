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
	DeliveryRepository *repository.DeliveryRepository
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountRepository, deliveryRepository *repository.DeliveryRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		OrderRepository:    orderRepository,
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
		AddressRepository:  addressRepository,
		DiscountRepository: discountRepository,
		DeliveryRepository: deliveryRepository,
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

	newProduct.Stock -= 1
	if err := c.ProductRepository.Update(tx, newProduct); err != nil {
		c.Log.Warnf("Failed to update stock of product : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newOrder.ProductId = newProduct.ID
	newOrder.ProductName = newProduct.Name
	newOrder.ProductDescription = newProduct.Description
	newOrder.Price = newProduct.Price
	newOrder.Quantity = request.Quantity
	newOrder.Amount = newProduct.Price * float32(newOrder.Quantity)

	// user/customer data
	newOrder.UserId = request.UserId
	newOrder.FirstName = request.FirstName
	newOrder.LastName = request.LastName
	newOrder.Email = request.Email
	newOrder.Phone = request.Phone

	// payment
	newOrder.PaymentMethod = request.PaymentMethod
	newOrder.PaymentStatus = helper.PENDING_PAYMENT

	if newOrder.PaymentMethod == helper.ONLINE {
		// jika pembayaran via online, fitur pengiriman "enabled"
		newOrder.IsDelivery = request.IsDelivery
		if newOrder.IsDelivery {
			newDelivery := new(entity.Delivery)
			if err := c.DeliveryRepository.FindFirst(tx, newDelivery); err != nil {
				c.Log.Warnf("Can't find delivery settings : %+v", err)
				return nil, fiber.ErrNotFound
			}
			newOrder.Distance = request.Distance
			newOrder.DeliveryCost = newOrder.Distance / newDelivery.Distance * newDelivery.Cost
			// jumlahkan semua total termasuk ongkir
			newOrder.Amount += newOrder.DeliveryCost

			// set status pengiriman
			newOrder.DeliveryStatus = helper.PREPARE_DELIVERY
		}
	} else if newOrder.PaymentMethod == helper.ONSITE {
		// jika pembayaran via onsite, fitur pengiriman "enabled"
		newOrder.DeliveryStatus = helper.PREPARE_DELIVERY
	}

	if request.DiscountCode != "" {
		newDiscount := new(entity.Discount)
		count, err := c.DiscountRepository.CountDiscountByCode(tx, newDiscount, request.DiscountCode)
		if err != nil {
			c.Log.Warnf("Failed to find discount by code : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// cek apakah diskonnya ada dan statusnya aktif (true)
		if count > 0 && newDiscount.Status {
			// cek apakah sudah kadaluarsa atau belum
			if newDiscount.End.After(time.Now()) {
				if newDiscount.Type == helper.PERCENT {
					newOrder.DiscountType = helper.PERCENT
					discount := float32(newDiscount.Value) / float32(100)
					afterDiscount := newOrder.Amount * discount
					newOrder.Amount -= afterDiscount
				} else if newDiscount.Type == helper.NOMINAL {
					newOrder.DiscountType = helper.NOMINAL
					newOrder.Amount -= newDiscount.Value
				}
			} else if newDiscount.End.Before(time.Now()) {
				c.Log.Warnf("Discount has expired : %+v", err)
				return nil, fiber.ErrBadRequest
			}
		} else if count > 0 && !newDiscount.Status {
			c.Log.Warnf("Discount has disabled : %+v", err)
			return nil, fiber.ErrBadRequest
		}
	} else if request.DiscountCode == "" {
		// default nominal
		newOrder.DiscountType = helper.NOMINAL
	}

	if err := c.ProductRepository.FindWithJoins(tx, newProduct, "Category"); err != nil {
		c.Log.Warnf("failed to find product with join category order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newOrder.CategoryName = newProduct.Category.Name
	newOrder.CompleteAddress = request.CompleteAddress
	newOrder.GoogleMapLink = request.GoogleMapLink

	if err := c.OrderRepository.Create(tx, newOrder); err != nil {
		c.Log.Warnf("failed to create new order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	invoice := fmt.Sprintf("INV/%d/CUST/%d", newOrder.ID, newOrder.UserId)
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

func (c *OrderUseCase) GetAllCurrent(ctx context.Context, request *model.GetOrderByCurrentRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("Failed to find all orders by user id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrdersToResponse(newOrders), nil
}

func (c *OrderUseCase) EditStatus(ctx context.Context, request *model.UpdateOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newOrder := new(entity.Order)
	newOrder.ID = request.ID
	if err := c.OrderRepository.FindById(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newOrder.PaymentStatus = request.PaymentStatus
	newOrder.DeliveryStatus = request.DeliveryStatus
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to update request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrderToResponse(newOrder), nil
}

func (c *OrderUseCase) GetByUserId(ctx context.Context, request *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("Failed to get all orders by user id from database : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrdersToResponse(newOrders), nil
}