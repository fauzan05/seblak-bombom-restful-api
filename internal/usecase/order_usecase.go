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
	"slices"
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
	WalletRepository       *repository.WalletRepository
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountCouponRepository, deliveryRepository *repository.DeliveryRepository,
	orderProductRepository *repository.OrderProductRepository, walletRepository *repository.WalletRepository) *OrderUseCase {
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
		WalletRepository:       walletRepository,
	}
}

func (c *OrderUseCase) Add(ctx context.Context, request *model.CreateOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newOrder := new(entity.Order)
	orderProducts := []entity.OrderProduct{}
	// temukan produk untuk memastikan ketersediaan dan masukkan data produk ke slice OrderProduct serta mengkalkulasikan tagihannya
	for _, orderProductRequest := range request.OrderProducts {
		if orderProductRequest.Quantity < 0 {
			c.Log.Warnf("Quantity must be positive number : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Quantity must be positive number!")
		}
		newProduct := new(entity.Product)
		newProduct.ID = orderProductRequest.ProductId
		count, err := c.ProductRepository.FindAndCountProductById(tx, newProduct)
		if count < 1 {
			c.Log.Warnf("Find product by id not found : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Product selected is not found!")
		}

		if newProduct.Stock < 1 {
			c.Log.Warnf("Product out of stock : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Product selected is out of stock!")
		}

		// pastikan permintaan tidak melebihi stok produk yang terkini
		newProduct.Stock -= orderProductRequest.Quantity
		if newProduct.Stock < 0 {
			c.Log.Warnf("Quantity order of product is out of limit : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Quantity order of product is out of limit")
		}

		// setelah dipastikan tidak melebihi stok produk yang terkini, kurangi stok produk terkini
		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			c.Log.Warnf("Failed to update stock of product : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "An error occurred on the server. Please try again later!")
		}

		orderProduct := entity.OrderProduct{
			ProductId:   orderProductRequest.ProductId,
			ProductName: newProduct.Name,
			Category:    newProduct.Category.Name,
			Price:       newProduct.Price,
			Quantity:    orderProductRequest.Quantity,
		}
		orderProducts = append(orderProducts, orderProduct)
		newOrder.Amount += orderProduct.Price * float32(orderProduct.Quantity)
	}

	newOrder.DeliveryCost = 0
	newOrder.IsDelivery = request.IsDelivery
	if newOrder.IsDelivery {
		// jika ingin dikirim, berarti ambil data delivery pada main address tiap user yang order
		newDelivery := new(entity.Delivery)
		newDelivery.ID = request.DeliveryId
		if err := c.DeliveryRepository.FindFirst(tx, newDelivery); err != nil {
			c.Log.Warnf("Can't find delivery settings : %+v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "Can't find delivery settings because not yet exist or not yet created")
		}

		// jumlahkan semua total termasuk ongkir
		newOrder.Amount += newDelivery.Cost
		newOrder.DeliveryCost = newDelivery.Cost
	}

	// user/customer data
	newOrder.UserId = request.UserId
	newOrder.FirstName = request.FirstName
	newOrder.LastName = request.LastName
	newOrder.Email = request.Email
	newOrder.Phone = request.Phone
	newOrder.Note = request.Note
	// set status order
	newOrder.OrderStatus = helper.ORDER_PENDING

	newOrder.DiscountType = helper.PERCENT
	if request.DiscountId > 0 {
		newDiscount := new(entity.DiscountCoupon)
		newDiscount.ID = request.DiscountId
		count, err := c.DiscountRepository.FindAndCountById(tx, newDiscount)
		if err != nil {
			c.Log.Warnf("Failed to find discount by code : %+v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Failed to find discount by code : %+v", err))
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
				c.Log.Warnf("Discount has expired!")
				return nil, fiber.NewError(fiber.StatusBadRequest, "Discount has expired!")
			}
		} else if count < 1 && !newDiscount.Status {
			c.Log.Warnf("Discount has disabled or doesn't exists!")
			return nil, fiber.NewError(fiber.StatusNotFound, "Discount has disabled or doesn't exists!")
		}
	}

	if !helper.IsValidPaymentMethod(request.PaymentMethod) {
		c.Log.Warnf("Invalid payment method!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid payment method!")
	}
	newOrder.PaymentMethod = request.PaymentMethod

	if !helper.IsValidChannelCode(request.ChannelCode) {
		c.Log.Warnf("Invalid channel code!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid channel code!")
	}
	newOrder.ChannelCode = request.ChannelCode

	if !helper.IsValidPaymentGateway(request.PaymentGateway) {
		c.Log.Warnf("Invalid payment gateway!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid payment gateway!")
	}
	newOrder.PaymentGateway = request.PaymentGateway

	if request.PaymentGateway == helper.PAYMENT_GATEWAY_SYSTEM {
		if request.PaymentMethod != helper.PAYMENT_METHOD_WALLET && request.ChannelCode != helper.WALLET_CHANNEL_CODE {
			c.Log.Warnf("Payment method %s is not available on payment gateway System!", request.PaymentMethod)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Payment method %s is not available on payment gateway System!", request.PaymentMethod))
		}
	}

	if request.PaymentGateway == helper.PAYMENT_GATEWAY_XENDIT {
		if request.PaymentMethod != helper.PAYMENT_METHOD_QR_CODE && request.PaymentMethod != helper.PAYMENT_METHOD_EWALLET {
			c.Log.Warnf("Payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway))
		} else {

			validChannelCodes := map[helper.PaymentMethod][]helper.ChannelCode{
				helper.PAYMENT_METHOD_QR_CODE: {
					helper.XENDIT_QR_DANA_CHANNEL_CODE,
					helper.XENDIT_QR_LINKAJA_CHANNEL_CODE,
				},
				helper.PAYMENT_METHOD_EWALLET: {
					helper.XENDIT_EWALLET_DANA_CHANNEL_CODE,
					helper.XENDIT_EWALLET_LINKAJA_CHANNEL_CODE,
					helper.XENDIT_EWALLET_OVO_CHANNEL_CODE,
					helper.XENDIT_EWALLET_SHOPEEPAY_CHANNEL_CODE,
				},
			}

			// Cek apakah ChannelCode valid untuk PaymentMethod yang dipilih
			isValid := false
			if validCodes, exists := validChannelCodes[request.PaymentMethod]; exists {
				if slices.Contains(validCodes, request.ChannelCode) {
					isValid = true
				}
			}

			// Jika tidak valid, berikan error
			if !isValid {
				c.Log.Warnf("Channel code %s is not available on payment gateway %s!", request.ChannelCode, request.PaymentGateway)
				return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Channel code %s is not available on payment gateway %s!", request.ChannelCode, request.PaymentGateway))
			}
		}
	}

	newOrder.PaymentStatus = helper.PENDING_PAYMENT
	if newOrder.PaymentMethod == helper.PAYMENT_METHOD_WALLET {
		// langsung paid dan proses walletnya
		if request.CurrentBalance < newOrder.Amount {
			// tampilkan error bahwa saldo kurang
			c.Log.Warnf("Your balance is insufficient to perform this transaction!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Your balance is insufficient to perform this transaction!")
		}

		newBalance := request.CurrentBalance - newOrder.Amount
		newWallet := new(entity.Wallet)
		if err := c.WalletRepository.UpdateWalletBalance(tx, newWallet, newOrder.UserId, newBalance); err != nil {
			c.Log.Warnf("Failed to update new balance : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to update new balance!")
		}

		newOrder.PaymentStatus = helper.PAID_PAYMENT
	}

	// mengambil alamat utama yang diambil oleh user
	newOrder.CompleteAddress = request.CompleteAddress

	if err := c.OrderRepository.Create(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to create new order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create new order : %+v", err))
	}

	timestamp := time.Now().Unix()
	dateStr := time.Now().Format("20060102")
	invoice := fmt.Sprintf("INV/%s/%d/ORDER/%d/CUST/%d", dateStr, timestamp, newOrder.ID, newOrder.UserId)
	newOrder.Invoice = invoice

	// memasukkan order_id ke order product
	for i := range orderProducts {
		orderProducts[i].OrderId = newOrder.ID
	}

	// insert semua data order product ke tabel order_products
	if err := c.OrderProductRepository.CreateInBatch(tx, &orderProducts); err != nil {
		c.Log.Warnf("Failed to add all order products into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to add all order products into database : %+v", err))
	}
	// mengisi kolom invoice ke tabel order setelah mendapatkan ID order nya
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to add invoice code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to add invoice code : %+v", err))
	}

	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find newly created order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find newly created order : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.OrderToResponse(newOrder), nil
}

func (c *OrderUseCase) GetAllCurrent(ctx context.Context, request *model.GetOrderByCurrentRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("Failed to find all orders by current user : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find all orders by user id : %+v", err))
	}

	return converter.OrdersToResponse(newOrders), nil
}

func (c *OrderUseCase) EditOrderStatus(ctx context.Context, request *model.UpdateOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newOrder := new(entity.Order)
	newOrder.ID = request.ID
	count, err := c.OrderRepository.FindAndCountById(tx, newOrder)
	if err != nil {
		c.Log.Warnf("Failed to find order by id into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find order by id into database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("Order not found by order id : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("Order not found by order id : %+v", err))
	}

	// validate first before update order status state into database
	// if rejected
	if request.OrderStatus == helper.ORDER_PENDING {
		if newOrder.OrderStatus == helper.ORDER_PENDING {
			c.Log.Warnf("Can't cancel an order that has been cancelled!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Can't cancel an order that has been cancelled!")
		}
		
		// find user wallet
		findWallet := new(entity.Wallet)
		if err := c.WalletRepository.FindEntityByUserId(tx, findWallet, newOrder.UserId); err != nil {
			c.Log.Warnf("Failed to find wallet by user id from database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find wallet by user id from database : %+v", err))
		}

		// return to wallet balance
		newWallet := new(entity.Wallet)
		newWallet.ID = findWallet.ID
		updateBalance := map[string]any{
			"balance": newOrder.Amount,
		}

		if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
			c.Log.Warnf("Failed to update wallet balance : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update wallet balance : %+v", err))
		}
	}

	newOrder.OrderStatus = request.OrderStatus
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("Failed to update status order by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to update status order by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.OrderToResponse(newOrder), nil
}

func (c *OrderUseCase) GetByUserId(ctx context.Context, request *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("Failed to get all orders by user id in the database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get all orders by user id from database : %+v", err))
	}

	return converter.OrdersToResponse(newOrders), nil
}
