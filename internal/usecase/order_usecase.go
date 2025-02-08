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
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	OrderRepository        *repository.OrderRepository
	ProductRepository      *repository.ProductRepository
	CategoryRepository     *repository.CategoryRepository
	AddressRepository      *repository.AddressRepository
	DiscountRepository     *repository.DiscountCouponRepository
	DeliveryRepository     *repository.DeliveryRepository
	OrderProductRepository *repository.OrderProductRepository
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountCouponRepository, deliveryRepository *repository.DeliveryRepository,
	orderProductRepository *repository.OrderProductRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		OrderRepository:        orderRepository,
		ProductRepository:      productRepository,
		CategoryRepository:     categoryRepository,
		AddressRepository:      addressRepository,
		DiscountRepository:     discountRepository,
		DeliveryRepository:     deliveryRepository,
		OrderProductRepository: orderProductRepository,
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
	orderProducts := []entity.OrderProduct{}
	// temukan produk untuk memastikan ketersediaan dan masukkan data produk ke slice OrderProduct serta mengkalkulasikan tagihannya
	for _, orderProductRequest := range request.OrderProducts {
		if orderProductRequest.Quantity < 0 {
			c.Log.Warnf("Quantity must be positive number : %+v", err)
			return nil, fiber.ErrBadRequest
		}
		newProduct := new(entity.Product)
		newProduct.ID = orderProductRequest.ProductId
		count, err := c.ProductRepository.FindAndCountById(tx, newProduct)
		if count < 1 {
			c.Log.Warnf("Find product by id not found : %+v", err)
			return nil, fiber.ErrNotFound
		}

		if newProduct.Stock < 1 {
			c.Log.Warnf("Product out of stock : %+v", err)
			return nil, fiber.ErrNotFound
		}

		// pastikan permintaan tidak melebihi stok produk yang terkini
		newProduct.Stock -= orderProductRequest.Quantity
		if newProduct.Stock < 0 {
			c.Log.Warnf("Quantity order of product is out of limit : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		// setelah dipastikan tidak melebihi stok produk yang terkini, kurangi stok produk terkini
		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			c.Log.Warnf("Failed to update stock of product : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		orderProduct := entity.OrderProduct{
			ProductId:   orderProductRequest.ProductId,
			ProductName: newProduct.Name,
			Price:       newProduct.Price,
			Quantity:    orderProductRequest.Quantity,
		}
		orderProducts = append(orderProducts, orderProduct)
		newOrder.Amount += orderProduct.Price * float32(orderProduct.Quantity)
	}

	// user/customer data
	newOrder.UserId = request.UserId
	newOrder.FirstName = request.FirstName
	newOrder.LastName = request.LastName
	newOrder.Email = request.Email
	newOrder.Phone = request.Phone
	newOrder.Note = request.Note

	// payment
	newOrder.PaymentMethod = request.PaymentMethod
	newOrder.PaymentStatus = helper.PENDING_PAYMENT

	newOrder.IsDelivery = request.IsDelivery
	if newOrder.IsDelivery {
		// jika ingin dikirim, berarti ambil data delivery pada main address tiap user yang order
		newDelivery := new(entity.Delivery)
		if err := c.DeliveryRepository.FindFirst(tx, newDelivery); err != nil {
			c.Log.Warnf("Can't find delivery settings : %+v", err)
			return nil, fiber.ErrNotFound
		}
		
		// jumlahkan semua total termasuk ongkir
		newOrder.Amount += newOrder.DeliveryCost
	} 

	// set status order
	newOrder.OrderStatus = helper.ORDER_PENDING

	if request.DiscountCode != "" {
		newDiscount := new(entity.DiscountCoupon)
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
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = afterDiscount
				} else if newDiscount.Type == helper.NOMINAL {
					newOrder.DiscountType = helper.NOMINAL
					newOrder.Amount -= newDiscount.Value
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = newDiscount.Value
				}
			} else if newDiscount.End.Before(time.Now()) {
				c.Log.Warnf("Discount has expired : %+v", err)
				return nil, fiber.ErrBadRequest
			}
		} else if count < 1 && !newDiscount.Status {
			c.Log.Warnf("Discount has disabled or doesn't exists : %+v", err)
			return nil, fiber.ErrBadRequest
		}
	} else if request.DiscountCode == "" {
		// default nominal
		newOrder.DiscountType = helper.NOMINAL
	}

	// mengambil alamat utama yang diambil oleh user
	newOrder.CompleteAddress = request.CompleteAddress

	if err := c.OrderRepository.Create(tx, newOrder); err != nil {
		c.Log.Warnf("failed to create new order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	invoice := fmt.Sprintf("INV/%d/CUST/%d", newOrder.ID, newOrder.UserId)
	newOrder.Invoice = invoice

	// memasukkan order_id ke order product
	for i := range orderProducts {
		orderProducts[i].OrderId = newOrder.ID
	}

	// insert semua data order product ke tabel order_products
	if err := c.OrderProductRepository.CreateInBatch(tx, &orderProducts); err != nil {
		c.Log.Warnf("failed to add all order products into database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	// mengisi kolom invoice ke tabel order setelah mendapatkan ID order nya
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("failed to add invoice code : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order with preload : %+v", err)
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
	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id into database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// cari apakah order by id itu ada di database
	if newOrder.Invoice == "" {
		c.Log.Warnf("Failed to find order by id (order not found) : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if newOrder.PaymentStatus == helper.FAILED_PAYMENT || newOrder.PaymentStatus == helper.PAID_PAYMENT {
		// jika ordernya statusnya ternyata sudah failed atau paid (berusaha untuk melakukan request ke 2x), maka tolak request tersebut agar stock produknya tidak ikut bertambah
		c.Log.Warnf("Failed to edit status order with has failed or paid payment status : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// jika ingin mengubah status menjadi gagal atau terbayar
	if request.PaymentStatus == helper.FAILED_PAYMENT {
		for _, orderProduct := range newOrder.OrderProducts {
			newProduct := new(entity.Product)
			newProduct.ID = orderProduct.ProductId
			// mencari data terkini dari produk dengan id
			if err := c.ProductRepository.FindById(tx, newProduct); err != nil {
				c.Log.Warnf("Failed to find product by id : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
			// tambahkan/kembalikan quantitas produk karena transaksinya gagal
			newProduct.Stock += orderProduct.Quantity
			// perbarui stok barang sekarang
			if err := c.ProductRepository.Update(tx, newProduct); err != nil {
				c.Log.Warnf("Failed to update product stock : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}
	}

	newOrder.PaymentStatus = request.PaymentStatus
	newOrder.OrderStatus = request.OrderStatus
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to update status order by id : %+v", err)
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
