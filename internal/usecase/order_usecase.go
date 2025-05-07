package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	xenditUseCase "seblak-bombom-restful-api/internal/usecase/xendit"
	"time"

	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB                             *gorm.DB
	Log                            *logrus.Logger
	Validate                       *validator.Validate
	OrderRepository                *repository.OrderRepository
	ProductRepository              *repository.ProductRepository
	CategoryRepository             *repository.CategoryRepository
	AddressRepository              *repository.AddressRepository
	DiscountRepository             *repository.DiscountCouponRepository
	DiscountUsageRepository        *repository.DiscountUsageRepository
	DeliveryRepository             *repository.DeliveryRepository
	OrderProductRepository         *repository.OrderProductRepository
	WalletRepository               *repository.WalletRepository
	XenditTransactionRepository    *repository.XenditTransctionRepository
	XenditTransactionQRCodeUseCase *xenditUseCase.XenditTransactionQRCodeUseCase
	XenditClient                   *xendit.APIClient
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountCouponRepository, discountUsageRepository *repository.DiscountUsageRepository,
	deliveryRepository *repository.DeliveryRepository, orderProductRepository *repository.OrderProductRepository,
	walletRepository *repository.WalletRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditTransactionQRCodeUseCase *xenditUseCase.XenditTransactionQRCodeUseCase, xenditClient *xendit.APIClient) *OrderUseCase {
	return &OrderUseCase{
		DB:                             db,
		Log:                            log,
		Validate:                       validate,
		OrderRepository:                orderRepository,
		ProductRepository:              productRepository,
		CategoryRepository:             categoryRepository,
		AddressRepository:              addressRepository,
		DiscountRepository:             discountRepository,
		DiscountUsageRepository:        discountUsageRepository,
		DeliveryRepository:             deliveryRepository,
		OrderProductRepository:         orderProductRepository,
		WalletRepository:               walletRepository,
		XenditTransactionRepository:    xenditTransactionRepository,
		XenditTransactionQRCodeUseCase: xenditTransactionQRCodeUseCase,
		XenditClient:                   xenditClient,
	}
}

func (c *OrderUseCase) Add(ctx *fiber.Ctx, request *model.CreateOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newOrder := new(entity.Order)
	orderProducts := []entity.OrderProduct{}
	var totalPriceOrderProduct float32
	// temukan produk untuk memastikan ketersediaan dan masukkan data produk ke slice OrderProduct serta mengkalkulasikan tagihannya
	for _, orderProductRequest := range request.OrderProducts {
		if orderProductRequest.Quantity < 0 {
			c.Log.Warnf("quantity must be positive number : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "quantity must be positive number!")
		}
		newProduct := new(entity.Product)
		newProduct.ID = orderProductRequest.ProductId
		count, err := c.ProductRepository.FindAndCountProductById(tx, newProduct)
		if count < 1 {
			c.Log.Warnf("product not found : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "product not found!")
		}

		if newProduct.Stock < 1 {
			c.Log.Warnf("product out of stock : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "product out of stock!")
		}

		// pastikan permintaan tidak melebihi stok produk yang terkini
		newProduct.Stock -= orderProductRequest.Quantity
		if newProduct.Stock < 0 {
			c.Log.Warnf("quantity order of product is out of limit : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "quantity order of product is out of limit")
		}

		// setelah dipastikan tidak melebihi stok produk yang terkini, kurangi stok produk terkini
		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			c.Log.Warnf("failed to update stock of product : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update stock of product : %+v", err))
		}

		orderProduct := entity.OrderProduct{
			ProductId:   orderProductRequest.ProductId,
			ProductName: newProduct.Name,
			Category:    newProduct.Category.Name,
			Price:       newProduct.Price,
			Quantity:    orderProductRequest.Quantity,
		}
		orderProducts = append(orderProducts, orderProduct)
		newOrder.TotalFinalPrice += orderProduct.Price * float32(orderProduct.Quantity)
		totalPriceOrderProduct = newOrder.TotalFinalPrice
	}

	newOrder.TotalProductPrice = newOrder.TotalFinalPrice
	newOrder.DeliveryCost = 0
	newOrder.IsDelivery = request.IsDelivery
	if newOrder.IsDelivery {
		// jika ingin dikirim, berarti ambil data delivery pada main address tiap user yang order
		newDelivery := new(entity.Delivery)
		newDelivery.ID = request.DeliveryId
		if err := c.DeliveryRepository.FindFirst(tx, newDelivery); err != nil {
			c.Log.Warnf("can't find delivery settings : %+v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("can't find delivery settings : %+v", err))
		}

		// jumlahkan semua total termasuk ongkir
		newOrder.TotalFinalPrice += newDelivery.Cost
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
			c.Log.Warnf("failed to find discount by code : %+v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to find discount by code : %+v", err))
		}

		// cek apakah diskonnya ada dan statusnya aktif (true)
		if count > 0 && newDiscount.Status {

			// cek apakah diskon masih berlaku pada waktu hari ini
			if newDiscount.End.After(time.Now()) && newDiscount.Start.Before(time.Now()) {
				// cek apakah minimal ordernya sudah sesuai
				if totalPriceOrderProduct < newDiscount.MinOrderValue {
					c.Log.Warnf("the order does not meet the minimum purchase requirements for this discount coupon!")
					return nil, fiber.NewError(fiber.StatusBadRequest, "the order does not meet the minimum purchase requirements for this discount coupon!")
				}

				// cek apakah user ini jatah diskonnya sudah habis atau belum
				discountUsage := new(entity.DiscountUsage)
				if err := c.DiscountUsageRepository.FindDiscountUsage(tx, discountUsage, newDiscount.ID, newOrder.UserId); err != nil {
					c.Log.Warnf("%+v", err)
					return nil, err
				}

				if discountUsage.ID > 0 {
					// maka update saja
					if discountUsage.UsageCount >= newDiscount.MaxUsagePerUser {
						c.Log.Warnf("the usage limit for this discount coupon has been exceeded!")
						return nil, fiber.NewError(fiber.StatusBadRequest, "the usage limit for this discount coupon has been exceeded!")
					}

					discountUsage.UsageCount = discountUsage.UsageCount + 1
					if err := c.DiscountUsageRepository.Update(tx, discountUsage); err != nil {
						c.Log.Warnf("failed to update usage count : %+v", err)
						return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update usage count : %+v", err))
					}
				} else {
					// maka create baru
					discountUsage.UsageCount = discountUsage.UsageCount + 1
					discountUsage.UserId = newOrder.UserId
					discountUsage.CouponId = newDiscount.ID
					if err := c.DiscountUsageRepository.Create(tx, discountUsage); err != nil {
						c.Log.Warnf("failed to create new usage count : %+v", err)
						return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create new usage count : %+v", err))
					}
				}

				if newDiscount.Type == helper.PERCENT {
					newOrder.DiscountType = helper.PERCENT
					discount := float32(newDiscount.Value) / float32(100)
					afterDiscount := newOrder.TotalFinalPrice * discount
					newOrder.TotalFinalPrice -= afterDiscount
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = afterDiscount
				} else if newDiscount.Type == helper.NOMINAL {
					newOrder.DiscountType = helper.NOMINAL
					newOrder.TotalFinalPrice -= newDiscount.Value
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = newDiscount.Value
				}

				newOrder.DiscountValue = newDiscount.Value
			} else {
				if newDiscount.End.Before(time.Now()) {
					c.Log.Warnf("discount has expired and is no longer available!")
					return nil, fiber.NewError(fiber.StatusBadRequest, "discount has expired and is no longer available!")
				} else if newDiscount.Start.After(time.Now()) {
					c.Log.Warnf("discount is not yet valid. It will be active starting %+s", newDiscount.Start.Format("January 02 2006 at 15:04:05"))
					return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("discount is not yet valid. It will be active starting %+s", newDiscount.Start.Format("January 02 2006 at 15:04:05")))
				}
			}
		} else if count < 1 && !newDiscount.Status {
			c.Log.Warnf("discount has disabled or doesn't exists!")
			return nil, fiber.NewError(fiber.StatusNotFound, "discount has disabled or doesn't exists!")
		}
	}

	if !helper.IsValidPaymentMethod(request.PaymentMethod) {
		c.Log.Warnf("invalid payment method!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid payment method!")
	}
	newOrder.PaymentMethod = request.PaymentMethod

	if !helper.IsValidChannelCode(request.ChannelCode) {
		c.Log.Warnf("invalid channel code!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid channel code!")
	}
	newOrder.ChannelCode = request.ChannelCode

	if !helper.IsValidPaymentGateway(request.PaymentGateway) {
		c.Log.Warnf("invalid payment gateway!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid payment gateway!")
	}
	newOrder.PaymentGateway = request.PaymentGateway
	newOrder.PaymentStatus = helper.PENDING_PAYMENT

	if request.PaymentGateway == helper.PAYMENT_GATEWAY_SYSTEM {
		if request.PaymentMethod != helper.PAYMENT_METHOD_WALLET {
			c.Log.Warnf("payment method %s is not available on payment gateway system!", request.PaymentMethod)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("payment method %s is not available on payment gateway system!", request.PaymentMethod))
		}

		if request.ChannelCode != helper.WALLET_CHANNEL_CODE {
			c.Log.Warnf("channel code %s is not available on payment gateway system!", request.ChannelCode)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("channel code %s is not available on payment gateway system!", request.ChannelCode))
		}

		// langsung paid dan proses walletnya
		if request.CurrentBalance < newOrder.TotalFinalPrice {
			// tampilkan error bahwa saldo kurang
			c.Log.Warnf("your balance is insufficient to perform this transaction!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "your balance is insufficient to perform this transaction!")
		}

		newBalance := request.CurrentBalance - newOrder.TotalFinalPrice
		newWallet := new(entity.Wallet)
		if err := c.WalletRepository.UpdateWalletBalance(tx, newWallet, newOrder.UserId, newBalance); err != nil {
			c.Log.Warnf("failed to update new balance : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "failed to update new balance!")
		}

		newOrder.PaymentStatus = helper.PAID_PAYMENT
	}

	if request.PaymentGateway == helper.PAYMENT_GATEWAY_XENDIT {
		if request.PaymentMethod != helper.PAYMENT_METHOD_QR_CODE {
			c.Log.Warnf("payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway))
		} else {

			validChannelCodes := map[helper.PaymentMethod][]helper.ChannelCode{
				helper.PAYMENT_METHOD_QR_CODE: {
					helper.XENDIT_QR_DANA_CHANNEL_CODE,
					helper.XENDIT_QR_LINKAJA_CHANNEL_CODE,
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
				c.Log.Warnf("channel code %s is not available on payment gateway %s!", request.ChannelCode, request.PaymentGateway)
				return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("channel code %s is not available on payment gateway %s!", request.ChannelCode, request.PaymentGateway))
			}
		}
	}

	// mengambil alamat utama yang diambil oleh user
	newOrder.CompleteAddress = request.CompleteAddress

	if err := c.OrderRepository.Create(tx, newOrder); err != nil {
		c.Log.Warnf("failed to create new order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create new order : %+v", err))
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
		c.Log.Warnf("failed to add all order products into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to add all order products into database : %+v", err))
	}

	// mengisi kolom invoice ke tabel order setelah mendapatkan ID order nya
	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("failed to add invoice code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to add invoice code : %+v", err))
	}

	// jika pembayaran menggunakan xendit,maka panggil xendit usecase
	if newOrder.PaymentGateway == helper.PAYMENT_GATEWAY_XENDIT && newOrder.PaymentMethod == helper.PAYMENT_METHOD_QR_CODE {
		newXenditQRCodeRequest := new(model.CreateXenditTransaction)
		newXenditQRCodeRequest.OrderId = newOrder.ID
		result, err := c.XenditTransactionQRCodeUseCase.Add(ctx, newXenditQRCodeRequest, tx)
		if err != nil {
			c.Log.Warn(err)
			return nil, err
		}

		newXenditTransaction := new(entity.XenditTransactions)
		newXenditTransaction.ID = result.ID
		newXenditTransaction.OrderId = result.OrderId
		newXenditTransaction.ReferenceId = result.ReferenceId
		newXenditTransaction.Amount = result.Amount
		newXenditTransaction.Currency = result.Currency
		newXenditTransaction.PaymentMethod = result.PaymentMethod
		newXenditTransaction.PaymentMethodId = result.PaymentMethodId
		newXenditTransaction.ChannelCode = result.ChannelCode
		newXenditTransaction.QrString = result.QrString
		newXenditTransaction.Status = result.Status
		newXenditTransaction.Description = result.Description
		newXenditTransaction.FailureCode = result.FailureCode
		newXenditTransaction.Metadata = result.Metadata
		newXenditTransaction.ExpiresAt = helper.TimeRFC3339.ToTime(result.ExpiresAt)
		newXenditTransaction.CreatedAt = helper.TimeRFC3339.ToTime(result.CreatedAt)
		newXenditTransaction.UpdatedAt = helper.TimeRFC3339.ToTime(result.UpdatedAt)

		newOrder.XenditTransaction = newXenditTransaction
	}

	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("failed to find newly created order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find newly created order : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.OrderToResponse(newOrder), nil
}

func (c *OrderUseCase) GetAllCurrent(ctx context.Context, request *model.GetOrderByCurrentRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("failed to find all orders by current user : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find all orders by user id : %+v", err))
	}

	return converter.OrdersToResponse(newOrders), nil
}

func (c *OrderUseCase) GetAllPaginate(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string, currentUser *model.UserResponse) (*[]model.OrderResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	if sortingColumn == "" {
		sortingColumn = "orders.id"
	}

	newPagination := new(repository.Pagination)
	newPagination.Page = page
	newPagination.PageSize = perPage
	newPagination.Column = sortingColumn
	newPagination.SortBy = sortBy
	allowedColumns := map[string]bool{
		"orders.id":                  true,
		"orders.invoice":             true,
		"orders.total_final_price":   true,
		"orders.total_product_price": true,
		"orders.discount_value":      true,
		"orders.discount_type":       true,
		"orders.total_discount":      true,
		"orders.user_id":             true,
		"orders.first_name":          true,
		"orders.last_name":           true,
		"orders.email":               true,
		"orders.phone":               true,
		"orders.payment_gateway":     true,
		"orders.payment_method":      true,
		"orders.channel_code":        true,
		"orders.payment_status":      true,
		"orders.order_status":        true,
		"orders.is_delivery":         true,
		"orders.delivery_cost":       true,
		"orders.complete_address":    true,
		"orders.note":                true,
		"orders.created_at":          true,
		"orders.updated_at":          true,
		"order_products.product_name": true,
		"order_products.category": true,
	}

	if !allowedColumns[newPagination.Column] {
		c.Log.Warnf("invalid sort column : %s", newPagination.Column)
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid sort column : %s", newPagination.Column))
	}

	orders, totalOrder, err := repository.Paginate(tx, &entity.Order{}, newPagination, func(d *gorm.DB) *gorm.DB {
		result := d.Joins("JOIN order_products ON order_products.order_id = orders.id").
			Preload("OrderProducts").
			Preload("OrderProducts.Product.Images").
			Preload("XenditTransaction").Where("order_products.product_name LIKE ?", "%"+search+"%")
		if currentUser.Role == helper.CUSTOMER {
			result.Where("user_id = ?", currentUser.ID)
		}
		return result
	})

	if err != nil {
		c.Log.Warnf("failed to paginate orders : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to paginate orders : %+v", err))
	}

	// Hitung total halaman
	var totalPages int = 0
	totalPages = int(totalOrder / int64(perPage))
	if totalOrder%int64(perPage) > 0 {
		totalPages++
	}

	return converter.OrdersToResponse(&orders), totalOrder, totalPages, nil
}

func (c *OrderUseCase) EditOrderStatus(ctx context.Context, request *model.UpdateOrderRequest, currentUser *model.UserResponse) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newOrder := new(entity.Order)
	newOrder.ID = request.ID
	count, err := c.OrderRepository.FindAndCountById(tx, newOrder)
	if err != nil {
		c.Log.Warnf("failed to find order by id into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find order by id into database : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("order not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "order not found!")
	}

	// validate first before update order status state into database
	// if cancelled
	if request.OrderStatus == helper.ORDER_CANCELLED {
		if newOrder.OrderStatus == helper.ORDER_CANCELLED || newOrder.OrderStatus == helper.ORDER_REJECTED || newOrder.OrderStatus == helper.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't cancel an order that has been rejected/cancelled/cancellation requested!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that has been rejected/cancelled/cancellation requested!")
		}

		if newOrder.OrderStatus == helper.ORDER_PENDING && newOrder.PaymentStatus == helper.PAID_PAYMENT {
			// jika pending dan paid maka kembalikan saldo
			// find user wallet
			findWallet := new(entity.Wallet)
			count, err := c.WalletRepository.FindAndCountFirstWalletByUserId(tx, findWallet, newOrder.UserId, "active")
			if err != nil {
				c.Log.Warnf("failed to find wallet by user id from database : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find wallet by user id from database : %+v", err))
			}

			if count < 1 {
				c.Log.Warnf("the selected wallet is not found!")
				return nil, fiber.NewError(fiber.StatusBadRequest, "the selected wallet is not found!")
			}

			// return to wallet balance
			newWallet := new(entity.Wallet)
			newWallet.ID = findWallet.ID
			totalBalance := newOrder.TotalFinalPrice + findWallet.Balance
			updateBalance := map[string]any{
				"balance": totalBalance,
			}

			if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
				c.Log.Warnf("failed to update wallet balance : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance : %+v", err))
			}

			newOrder.OrderStatus = request.OrderStatus
		} else if newOrder.OrderStatus == helper.ORDER_PENDING && newOrder.PaymentStatus != helper.PENDING_PAYMENT {
			newOrder.PaymentStatus = helper.CANCELLED_PAYMENT
			newOrder.OrderStatus = request.OrderStatus
		} else if newOrder.OrderStatus == helper.ORDER_RECEIVED {
			// memerlukan persetujuan seller
			newOrder.OrderStatus = helper.ORDER_CANCELLATION_REQUESTED
		} else if newOrder.OrderStatus == helper.READY_FOR_PICKUP {
			c.Log.Warnf("can't cancel an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that is ready for pickup!")
		} else if newOrder.OrderStatus == helper.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't cancel an order that is being delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that is being delivered!")
		} else if newOrder.OrderStatus == helper.ORDER_DELIVERED {
			c.Log.Warnf("can't cancel an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that has been delivered!")
		}

	} 
	
	if request.OrderStatus == helper.ORDER_REJECTED {
		// Admin access only for reject
		if currentUser.Role == helper.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.OrderStatus == helper.ORDER_REJECTED {
			c.Log.Warnf("can't reject an order that has been rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been rejected!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLED {
			c.Log.Warnf("can't reject an order that has been cancelled!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been cancelled!")
		}

		if newOrder.OrderStatus == helper.ORDER_RECEIVED {
			c.Log.Warnf("can't reject an order that has been received!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been received!")
		}

		if newOrder.OrderStatus == helper.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't reject an order that is been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that is been delivered!")
		}

		// maka balikkan saldo customer
		if newOrder.PaymentStatus == helper.PAID_PAYMENT {
			// find user wallet
			// kembalikan saldo ke user
			findWallet := new(entity.Wallet)
			count, err := c.WalletRepository.FindAndCountFirstWalletByUserId(tx, findWallet, newOrder.UserId, "active")
			if err != nil {
				c.Log.Warnf("failed to find wallet by user id from database : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find wallet by user id from database : %+v", err))
			}

			if count < 1 {
				c.Log.Warnf("the selected wallet is not found!")
				return nil, fiber.NewError(fiber.StatusBadRequest, "the selected wallet is not found!")
			}

			// return to wallet balance
			newWallet := new(entity.Wallet)
			newWallet.ID = findWallet.ID
			totalBalance := newOrder.TotalFinalPrice + findWallet.Balance
			updateBalance := map[string]any{
				"balance": totalBalance,
			}

			if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
				c.Log.Warnf("failed to update wallet balance : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance : %+v", err))
			}
		}

		newOrder.OrderStatus = request.OrderStatus
	}
	
	if request.OrderStatus == helper.ORDER_RECEIVED {
		// Admin access only for received
		if currentUser.Role == helper.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.PaymentStatus != helper.PAID_PAYMENT {
			c.Log.Warnf("can't accept an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLED || newOrder.OrderStatus == helper.ORDER_REJECTED {
			c.Log.Warnf("can't accept an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == helper.ORDER_RECEIVED {
			c.Log.Warnf("can't accept an order that has been received!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been received!")
		}

		if newOrder.OrderStatus == helper.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't accept an order that is being delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that is being delivered!")
		}

		if newOrder.OrderStatus == helper.ORDER_DELIVERED {
			c.Log.Warnf("can't accept an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been delivered!")
		}

		if newOrder.OrderStatus == helper.READY_FOR_PICKUP {
			c.Log.Warnf("can't accept an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that is ready for pickup!")
		}

		newOrder.OrderStatus = request.OrderStatus
	}
	
	if request.OrderStatus == helper.READY_FOR_PICKUP {
		// Admin access only for pick up
		if currentUser.Role == helper.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.PaymentStatus != helper.PAID_PAYMENT {
			c.Log.Warnf("can't pick up an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't pick up an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that is ready for pickup!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLED || newOrder.OrderStatus == helper.ORDER_REJECTED {
			c.Log.Warnf("can't pick up an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't pick up an order that has a cancellation request!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has a cancellation request!")
		}

		if newOrder.OrderStatus == helper.ORDER_DELIVERED {
			c.Log.Warnf("can't pick up an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has been delivered!")
		}

		newOrder.OrderStatus = request.OrderStatus
	}
	
	if request.OrderStatus == helper.ORDER_DELIVERED {
		if newOrder.PaymentStatus != helper.PAID_PAYMENT {
			c.Log.Warnf("can't complete an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == helper.ORDER_DELIVERED {
			c.Log.Warnf("can't complete an order that has been completed!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has been completed!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLED || newOrder.OrderStatus == helper.ORDER_REJECTED {
			c.Log.Warnf("can't complete an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == helper.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't complete an order that has a cancellation request!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has a cancellation request!")
		}

		newOrder.OrderStatus = request.OrderStatus
	}

	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("failed to update status order by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update status order by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.OrderToResponse(newOrder), nil
}

func (c *OrderUseCase) GetByUserId(ctx context.Context, request *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newOrders := new([]entity.Order)
	if err := c.OrderRepository.FindAllOrdersByUserId(tx, newOrders, request.ID); err != nil {
		c.Log.Warnf("failed to get all orders by user id in the database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get all orders by user id from database : %+v", err))
	}

	return converter.OrdersToResponse(newOrders), nil
}
