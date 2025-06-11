package usecase

import (
	"context"
	"fmt"
	"html/template"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	xenditUseCase "seblak-bombom-restful-api/internal/usecase/xendit"
	"strings"
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
	ApplicationRepository          *repository.ApplicationRepository
	NotificationRepository         *repository.NotificationRepository
	Email                          *mailer.EmailWorker
}

func NewOrderUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository,
	categoryRepository *repository.CategoryRepository, addressRepository *repository.AddressRepository,
	discountRepository *repository.DiscountCouponRepository, discountUsageRepository *repository.DiscountUsageRepository,
	deliveryRepository *repository.DeliveryRepository, orderProductRepository *repository.OrderProductRepository,
	walletRepository *repository.WalletRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditTransactionQRCodeUseCase *xenditUseCase.XenditTransactionQRCodeUseCase, xenditClient *xendit.APIClient,
	applicationRepository *repository.ApplicationRepository, email *mailer.EmailWorker, notificationRepository *repository.NotificationRepository) *OrderUseCase {
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
		ApplicationRepository:          applicationRepository,
		Email:                          email,
		NotificationRepository:         notificationRepository,
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
	productsSelected := []map[string]any{}
	var totalPriceOrderProduct float32
	// temukan produk untuk memastikan ketersediaan dan masukkan data produk ke slice OrderProduct serta mengkalkulasikan tagihannya
	for _, orderProductRequest := range request.OrderProducts {
		if orderProductRequest.Quantity < 0 {
			c.Log.Warnf("quantity must be positive number!")
			return nil, fiber.NewError(fiber.StatusInternalServerError, "quantity must be positive number!")
		}
		newProduct := new(entity.Product)
		newProduct.ID = orderProductRequest.ProductId
		count, err := c.ProductRepository.FindAndCountProductById(tx, newProduct)
		if err != nil {
			c.Log.Warnf("failed to find product by id : %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to find product by id : %+v", err))
		}

		if count < 1 {
			c.Log.Warnf("product not found!")
			return nil, fiber.NewError(fiber.StatusInternalServerError, "product not found!")
		}

		var imageSelectedFileName string
		var productImageBase64 string
		for _, image := range newProduct.Images {
			if image.Position == 1 {
				imageSelectedFileName = image.FileName
				productImagePath := fmt.Sprintf("../uploads/images/products/%s", image.FileName)
				imageBase64, err := helper_others.ImageToBase64(productImagePath)
				if err != nil {
					c.Log.Warnf("failed to convert product image to base64 : %+v", err)
					return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert product image to base64 : %+v", err))
				}
				productImageBase64 = imageBase64
			}
			break
		}

		productSelected := map[string]any{
			"ProductImageFilename": imageSelectedFileName,
			"ProductImage":         productImageBase64,
			"ProductName":          newProduct.Name,
			"Quantity":             orderProductRequest.Quantity,
			"Price":                helper_others.FormatNumberFloat32(newProduct.Price),
		}

		productsSelected = append(productsSelected, productSelected)

		if newProduct.Stock < 1 {
			c.Log.Warnf("product out of stock!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "product out of stock!")
		}

		// pastikan permintaan tidak melebihi stok produk yang terkini
		newProduct.Stock -= orderProductRequest.Quantity
		if newProduct.Stock < 0 {
			c.Log.Warnf("quantity order of product is out of limit!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "quantity order of product is out of limit")
		}

		// setelah dipastikan tidak melebihi stok produk yang terkini, kurangi stok produk terkini
		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			c.Log.Warnf("failed to update stock of product : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update stock of product : %+v", err))
		}

		orderProduct := entity.OrderProduct{
			ProductId:                 orderProductRequest.ProductId,
			ProductName:               newProduct.Name,
			ProductFirstImagePosition: imageSelectedFileName,
			Category:                  newProduct.Category.Name,
			Price:                     newProduct.Price,
			Quantity:                  orderProductRequest.Quantity,
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
	newOrder.OrderStatus = enum_state.ORDER_PENDING

	newOrder.DiscountType = enum_state.PERCENT
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

				if newDiscount.Type == enum_state.PERCENT {
					newOrder.DiscountType = enum_state.PERCENT
					discount := float32(newDiscount.Value) / float32(100)
					afterDiscount := newOrder.TotalFinalPrice * discount
					newOrder.TotalFinalPrice -= afterDiscount
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = afterDiscount
				} else if newDiscount.Type == enum_state.NOMINAL {
					newOrder.DiscountType = enum_state.NOMINAL
					newOrder.TotalFinalPrice -= newDiscount.Value
					// simpan total diskon/potongan harganya
					newOrder.TotalDiscount = newDiscount.Value
				}

				newOrder.DiscountValue = newDiscount.Value
				// update used count
				newDiscount.UsedCount = newDiscount.UsedCount + 1
				if err := c.DiscountRepository.Update(tx, newDiscount); err != nil {
					c.Log.Warnf("failed to update discount used count : %+v", err)
					return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update discount used count : %+v", err))
				}
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

	if !enum_state.IsValidPaymentMethod(request.PaymentMethod) {
		c.Log.Warnf("invalid payment method!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid payment method!")
	}
	newOrder.PaymentMethod = request.PaymentMethod

	if !enum_state.IsValidChannelCode(request.ChannelCode) {
		c.Log.Warnf("invalid channel code!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid channel code!")
	}
	newOrder.ChannelCode = request.ChannelCode

	if !enum_state.IsValidPaymentGateway(request.PaymentGateway) {
		c.Log.Warnf("invalid payment gateway!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid payment gateway!")
	}
	newOrder.PaymentGateway = request.PaymentGateway
	newOrder.PaymentStatus = enum_state.PENDING_PAYMENT

	if request.PaymentGateway == enum_state.PAYMENT_GATEWAY_SYSTEM {
		if request.PaymentMethod != enum_state.PAYMENT_METHOD_WALLET {
			c.Log.Warnf("payment method %s is not available on payment gateway system!", request.PaymentMethod)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("payment method %s is not available on payment gateway system!", request.PaymentMethod))
		}

		if request.ChannelCode != enum_state.WALLET_CHANNEL_CODE {
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

		newOrder.PaymentStatus = enum_state.PAID_PAYMENT
	}

	if request.PaymentGateway == enum_state.PAYMENT_GATEWAY_XENDIT {
		if request.PaymentMethod != enum_state.PAYMENT_METHOD_QR_CODE {
			c.Log.Warnf("payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway)
			return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("payment method %s is not available on payment gateway %s!", request.PaymentMethod, request.PaymentGateway))
		} else {

			validChannelCodes := map[enum_state.PaymentMethod][]enum_state.ChannelCode{
				enum_state.PAYMENT_METHOD_QR_CODE: {
					enum_state.XENDIT_QR_DANA_CHANNEL_CODE,
					enum_state.XENDIT_QR_LINKAJA_CHANNEL_CODE,
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

	newApp := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
		c.Log.Warnf("failed to find application from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
	}

	newOrder.ServiceFee = newApp.ServiceFee
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

	if request.PaymentMethod != enum_state.PAYMENT_METHOD_WALLET {
		now := time.Now()
		err = helper_others.SaveWalletTransaction(tx, request.UserId, &newOrder.ID, newOrder.TotalFinalPrice, enum_state.WALLET_FLOW_TYPE_DEBIT, enum_state.WALLET_TRANSACTION_TYPE_ORDER_PAYMENT,
			request.PaymentMethod, enum_state.WALLET_TRANSACTION_STATUS_COMPLETED, "", invoice, "", nil, &now)
		if err != nil {
			c.Log.Warnf("failed to save wallet transaction : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
		}
	}

	// jika pembayaran menggunakan xendit,maka panggil xendit usecase
	if newOrder.PaymentGateway == enum_state.PAYMENT_GATEWAY_XENDIT && newOrder.PaymentMethod == enum_state.PAYMENT_METHOD_QR_CODE {
		newXenditQRCodeRequest := new(model.CreateXenditTransaction)
		newXenditQRCodeRequest.OrderId = newOrder.ID
		newXenditQRCodeRequest.Lang = request.Lang
		newXenditQRCodeRequest.TimeZone = request.TimeZone
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
		newXenditTransaction.ExpiresAt = helper_others.TimeRFC3339.ToTime(result.ExpiresAt)
		newXenditTransaction.CreatedAt = helper_others.TimeRFC3339.ToTime(result.CreatedAt)
		newXenditTransaction.UpdatedAt = helper_others.TimeRFC3339.ToTime(result.UpdatedAt)

		newOrder.XenditTransaction = newXenditTransaction
	}

	// tidak perlu preload xendit_transactions karena sudah di handle pada if diatas
	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("failed to find newly created order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find newly created order : %+v", err))
	}

	if newApp.LogoFilename == "" {
		c.Log.Warnf("application logo has not uploaded yet!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "application logo has not uploaded yet!")
	}

	logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
	logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
	if err != nil {
		c.Log.Warnf("failed to convert logo to base64 : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
	}

	if newOrder.PaymentStatus == enum_state.PAID_PAYMENT {
		newMail := new(model.Mail)
		newMail.To = []string{newOrder.Email}
		newMail.Subject = "Payment Successfull"
		if request.Lang == enum_state.INDONESIA {
			newMail.Subject = "Pembayaran Berhasil"
		}

		baseTemplatePath := "../internal/templates/base_template_email1.html"
		childPath := fmt.Sprintf("../internal/templates/%s/email/order_payment.html", request.Lang)
		tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		paymentLink := fmt.Sprintf("%s/payment/orders/%d/details", request.BaseFrontEndURL, newOrder.ID)
		orderTrackingURL := fmt.Sprintf("%s/orders/%d/details", request.BaseFrontEndURL, newOrder.ID)
		bodyBuilder := new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]any{
			"CustomerName":     newOrder.FirstName + " " + newOrder.LastName,
			"Invoice":          newOrder.Invoice,
			"Date":             newOrder.CreatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
			"PaymentMethod":    string(newOrder.PaymentMethod),
			"Items":            productsSelected,
			"LogoImage":        logoImageBase64,
			"CompanyTitle":     newApp.AppName,
			"TotalAmount":      helper_others.FormatNumberFloat32(newOrder.TotalFinalPrice),
			"Year":             time.Now().Format("2006"),
			"CustomerNotes":    newOrder.Note,
			"ShippingMethod":   newOrder.IsDelivery,
			"ShippingCost":     helper_others.FormatNumberFloat32(newOrder.DeliveryCost),
			"ServiceFee":       helper_others.FormatNumberFloat32(newOrder.ServiceFee),
			"Discount":         helper_others.FormatNumberFloat32(newOrder.TotalDiscount),
			"Subject":          newMail.Subject,
			"PaymentStatus":    newOrder.PaymentStatus,
			"PaymentLink":      paymentLink,
			"OrderTrackingURL": orderTrackingURL,
		})
		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}

		newMail.Template = *bodyBuilder
		c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
		// send email
		select {
		case c.Email.MailQueue <- *newMail:
		default:
			c.Log.Warnf("email queue full, failed to send to %s", newOrder.Email)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newOrder.Email))
		}

		newNotification := new(entity.Notification)
		newNotification.UserID = newOrder.UserId
		newNotification.Title = newMail.Subject
		newNotification.IsRead = false
		newNotification.Type = enum_state.TRANSACTION
		baseTemplatePath = "../internal/templates/base_template_notification1.html"
		childPath = fmt.Sprintf("../internal/templates/%s/notification/order_payment.html", request.Lang)
		tmpl, err = template.ParseFiles(baseTemplatePath, childPath)
		if err != nil {
			c.Log.Warnf("failed to parse template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
		}

		logoImagePath = fmt.Sprintf("%s://%s/api/image/application/%s", ctx.Protocol(), ctx.Hostname(), newApp.LogoFilename)
		bodyBuilder = new(strings.Builder)
		err = tmpl.ExecuteTemplate(bodyBuilder, "base", map[string]string{
			"FirstName":        newOrder.FirstName,
			"Year":             time.Now().Format("2006"),
			"CompanyName":      newApp.AppName,
			"LogoImagePath":    logoImagePath,
			"Date":             newOrder.CreatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
			"Invoice":          newOrder.Invoice,
			"PaymentMethod":    string(newOrder.PaymentMethod),
			"Subject":          newMail.Subject,
			"PaymentStatus":    string(newOrder.PaymentStatus),
			"PaymentLink":      paymentLink,
			"OrderTrackingURL": orderTrackingURL,
		})

		if err != nil {
			c.Log.Warnf("failed to execute template file html : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
		}

		newNotification.BodyContent = bodyBuilder.String()
		if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
			c.Log.Warnf("failed to create notification into database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
		}
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
		"orders.id":                   true,
		"orders.invoice":              true,
		"orders.total_final_price":    true,
		"orders.total_product_price":  true,
		"orders.discount_value":       true,
		"orders.discount_type":        true,
		"orders.total_discount":       true,
		"orders.user_id":              true,
		"orders.first_name":           true,
		"orders.last_name":            true,
		"orders.email":                true,
		"orders.phone":                true,
		"orders.payment_gateway":      true,
		"orders.payment_method":       true,
		"orders.channel_code":         true,
		"orders.payment_status":       true,
		"orders.order_status":         true,
		"orders.is_delivery":          true,
		"orders.delivery_cost":        true,
		"orders.complete_address":     true,
		"orders.note":                 true,
		"orders.created_at":           true,
		"orders.updated_at":           true,
		"order_products.product_name": true,
		"order_products.category":     true,
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
		if currentUser.Role == enum_state.CUSTOMER {
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

	if request.OrderStatus == enum_state.ORDER_CANCELLED {
		newOrder.CancellationNotes = request.CancellationNotes
	}

	if request.OrderStatus == enum_state.ORDER_RECEIVED {
		newOrder.RejectionNotes = request.RejectionNotes
	}

	var is_send_email bool
	var mail_subject_cust string
	var mail_subject_admin string
	// validate first before update order status state into database
	// if cancelled
	if request.OrderStatus == enum_state.ORDER_CANCELLED {
		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED || newOrder.OrderStatus == enum_state.ORDER_REJECTED || newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't cancel an order that has been rejected/cancelled/cancellation requested!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that has been rejected/cancelled/cancellation requested!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_PENDING && newOrder.PaymentStatus == enum_state.PAID_PAYMENT {
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

			err = helper_others.SaveWalletTransaction(tx, findWallet.UserId, &newOrder.ID, newOrder.TotalFinalPrice,
				enum_state.WALLET_FLOW_TYPE_CREDIT, enum_state.WALLET_TRANSACTION_TYPE_ORDER_REFUND, newOrder.PaymentMethod,
				enum_state.WALLET_TRANSACTION_STATUS_COMPLETED, "", request.CancellationNotes, request.CancellationNotes, nil, nil)

			if err != nil {
				c.Log.Warnf("failed to save wallet transaction : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
			}

			newOrder.OrderStatus = request.OrderStatus
			is_send_email = true
			mail_subject_cust = fmt.Sprintf("Your Order with ID %d Has Been Cancelled", newOrder.ID)
			mail_subject_admin = fmt.Sprintf("Order ID %d Has Been Canceled by Customer", newOrder.ID)
			if request.Lang == enum_state.INDONESIA {
				mail_subject_cust = fmt.Sprintf("Pesanan Anda dengan ID %d Telah Dibatalkan", newOrder.ID)
				mail_subject_admin = fmt.Sprintf("Order ID %d Telah Dibatalkan oleh Customer", newOrder.ID)
			}
		} else if newOrder.OrderStatus == enum_state.ORDER_PENDING && newOrder.PaymentStatus != enum_state.PENDING_PAYMENT {
			newOrder.PaymentStatus = enum_state.CANCELLED_PAYMENT
			newOrder.OrderStatus = request.OrderStatus
		} else if newOrder.OrderStatus == enum_state.ORDER_RECEIVED {
			c.Log.Warnf("can't cancel an order that is received!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that is received!")
		} else if newOrder.OrderStatus == enum_state.READY_FOR_PICKUP {
			c.Log.Warnf("can't cancel an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that is ready for pickup!")
		} else if newOrder.OrderStatus == enum_state.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't cancel an order that is being delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that is being delivered!")
		} else if newOrder.OrderStatus == enum_state.ORDER_DELIVERED {
			c.Log.Warnf("can't cancel an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't cancel an order that has been delivered!")
		}
	}

	if request.OrderStatus == enum_state.ORDER_REJECTED {
		// Admin access only for reject
		if currentUser.Role == enum_state.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_REJECTED {
			c.Log.Warnf("can't reject an order that has been rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been rejected!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED {
			c.Log.Warnf("can't reject an order that has been cancelled!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been cancelled!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_RECEIVED {
			c.Log.Warnf("can't reject an order that has been received!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that has been received!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't reject an order that is been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't reject an order that is been delivered!")
		}

		// maka balikkan saldo customer
		if newOrder.PaymentStatus == enum_state.PAID_PAYMENT {
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

			now := time.Now()
			err = helper_others.SaveWalletTransaction(tx, findWallet.UserId, &newOrder.ID, newOrder.TotalFinalPrice,
				enum_state.WALLET_FLOW_TYPE_CREDIT, enum_state.WALLET_TRANSACTION_TYPE_ORDER_REFUND, newOrder.PaymentMethod,
				enum_state.WALLET_TRANSACTION_STATUS_COMPLETED, "", "", request.RejectionNotes, &currentUser.ID, &now)

			if err != nil {
				c.Log.Warnf("failed to save wallet transaction : %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
			}

			is_send_email = true
			mail_subject_cust = fmt.Sprintf("Your Order with ID %d Has Been Rejected", newOrder.ID)
			mail_subject_admin = fmt.Sprintf("Order ID %d Has Been Rejected by Admin", newOrder.ID)
			if request.Lang == enum_state.INDONESIA {
				mail_subject_cust = fmt.Sprintf("Pesanan Anda dengan ID %d Telah Ditolak", newOrder.ID)
				mail_subject_admin = fmt.Sprintf("Order ID %d Telah Ditolak oleh Admin", newOrder.ID)
			}
		}

		newOrder.OrderStatus = request.OrderStatus
	}

	if request.OrderStatus == enum_state.ORDER_RECEIVED {
		// Admin access only for received
		if currentUser.Role == enum_state.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.PaymentStatus != enum_state.PAID_PAYMENT {
			c.Log.Warnf("can't accept an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED || newOrder.OrderStatus == enum_state.ORDER_REJECTED {
			c.Log.Warnf("can't accept an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_RECEIVED {
			c.Log.Warnf("can't accept an order that has been received!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been received!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't accept an order that is being delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that is being delivered!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_DELIVERED {
			c.Log.Warnf("can't accept an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that has been delivered!")
		}

		if newOrder.OrderStatus == enum_state.READY_FOR_PICKUP {
			c.Log.Warnf("can't accept an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't accept an order that is ready for pickup!")
		}

		newOrder.OrderStatus = request.OrderStatus
		is_send_email = true
		mail_subject_cust = fmt.Sprintf("Your Order with ID %d Has Been Received", newOrder.ID)
		mail_subject_admin = fmt.Sprintf("Order ID %d Has Been Received by Admin", newOrder.ID)
		if request.Lang == enum_state.INDONESIA {
			mail_subject_cust = fmt.Sprintf("Pesanan Anda dengan ID %d Telah Diterima", newOrder.ID)
			mail_subject_admin = fmt.Sprintf("Order ID %d Telah Diterima oleh Admin", newOrder.ID)
		}
	}

	if request.OrderStatus == enum_state.READY_FOR_PICKUP {
		// Admin access only for pick up
		if currentUser.Role == enum_state.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.PaymentStatus != enum_state.PAID_PAYMENT {
			c.Log.Warnf("can't pick up an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't pick up an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that is ready for pickup!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED || newOrder.OrderStatus == enum_state.ORDER_REJECTED {
			c.Log.Warnf("can't pick up an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't pick up an order that has a cancellation request!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has a cancellation request!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_BEING_DELIVERED {
			c.Log.Warnf("can't pick up an order that has being delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has being delivered!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_DELIVERED {
			c.Log.Warnf("can't pick up an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't pick up an order that has been delivered!")
		}

		newOrder.OrderStatus = request.OrderStatus
	}

	if request.OrderStatus == enum_state.ORDER_BEING_DELIVERED {
		// Admin access only for being delivered
		if currentUser.Role == enum_state.CUSTOMER {
			c.Log.Warn("admin access only!")
			return nil, fiber.NewError(fiber.StatusUnauthorized, "admin access only!")
		}

		if newOrder.PaymentStatus != enum_state.PAID_PAYMENT {
			c.Log.Warnf("can't being delivered an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't being delivered an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't being delivered an order that is ready for pickup!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't being delivered an order that is ready for pickup!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED || newOrder.OrderStatus == enum_state.ORDER_REJECTED {
			c.Log.Warnf("can't being delivered an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't being delivered an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't being delivered an order that has a cancellation request!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't being delivered an order that has a cancellation request!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_DELIVERED {
			c.Log.Warnf("can't being delivered an order that has been delivered!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't being delivered an order that has been delivered!")
		}

		newOrder.OrderStatus = request.OrderStatus
	}

	if request.OrderStatus == enum_state.ORDER_DELIVERED {
		if newOrder.PaymentStatus != enum_state.PAID_PAYMENT {
			c.Log.Warnf("can't complete an order that has not been paid yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has not been paid yet!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_DELIVERED {
			c.Log.Warnf("can't complete an order that has been completed!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has been completed!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLED || newOrder.OrderStatus == enum_state.ORDER_REJECTED {
			c.Log.Warnf("can't complete an order that has been cancelled/rejected!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has been cancelled/rejected!")
		}

		if newOrder.OrderStatus == enum_state.ORDER_CANCELLATION_REQUESTED {
			c.Log.Warnf("can't complete an order that has a cancellation request!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "can't complete an order that has a cancellation request!")
		}

		newOrder.OrderStatus = request.OrderStatus
		is_send_email = true
		mail_subject_cust = fmt.Sprintf("Your Order with ID %d Has Been Picked Up", newOrder.ID)
		mail_subject_admin = fmt.Sprintf("Order ID %d Has Been Picked Up by Customer", newOrder.ID)
		if request.Lang == enum_state.INDONESIA {
			mail_subject_cust = fmt.Sprintf("Pesanan Anda dengan ID %d Telah Diambil", newOrder.ID)
			mail_subject_admin = fmt.Sprintf("Pesanan dengan ID %d Telah Diambil oleh Pelanggan", newOrder.ID)

		}
		if newOrder.IsDelivery {
			mail_subject_cust = fmt.Sprintf("Your Order with ID %d Has Been Delivered", newOrder.ID)
			mail_subject_admin = fmt.Sprintf("Order ID %d Has Been Delivered", newOrder.ID)
			if request.Lang == enum_state.INDONESIA {
				mail_subject_cust = fmt.Sprintf("Pesanan Anda dengan ID %d Telah Sampai Di Tujuan", newOrder.ID)
				mail_subject_admin = fmt.Sprintf("Order ID %d Telah Sampai Di Tujuan", newOrder.ID)
			}
		}
	}

	if err := c.OrderRepository.Update(tx, newOrder); err != nil {
		c.Log.Warnf("failed to update status order by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update status order by id : %+v", err))
	}

	if err := c.OrderRepository.FindWithPreloads(tx, newOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("failed to find newly created order : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find newly created order : %+v", err))
	}

	if is_send_email {
		// kirim email ke customer
		newApp := new(entity.Application)
		if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
			c.Log.Warnf("failed to find application from database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
		}

		if newApp.LogoFilename == "" {
			c.Log.Warnf("application logo has not uploaded yet!")
			return nil, fiber.NewError(fiber.StatusBadRequest, "application logo has not uploaded yet!")
		}

		logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
		logoImageBase64, err := helper_others.ImageToBase64(logoImagePath)
		if err != nil {
			c.Log.Warnf("failed to convert logo to base64 : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
		}

		baseTemplatePath := "../internal/templates/base_template_email1.html"
		childPath := fmt.Sprintf("../internal/templates/%s/email/order_status_customer.html", request.Lang)
		orderTrackingURL := fmt.Sprintf("%s/orders/%d/details", request.BaseFrontEndURL, newOrder.ID)
		data := map[string]any{
			"CustomerName":      newOrder.FirstName + " " + newOrder.LastName,
			"CustomerPhone":     newOrder.Phone,
			"Invoice":           newOrder.Invoice,
			"Date":              newOrder.UpdatedAt.In(&request.TimeZone).Format("02 Jan 2006 15:04 MST"),
			"PaymentMethod":     string(newOrder.PaymentMethod),
			"LogoImage":         logoImageBase64,
			"CompanyTitle":      newApp.AppName,
			"TotalAmount":       helper_others.FormatNumberFloat32(newOrder.TotalFinalPrice),
			"Year":              time.Now().Format("2006"),
			"PaymentStatus":     newOrder.PaymentStatus,
			"OrderTrackingURL":  orderTrackingURL,
			"CancellationNotes": newOrder.CancellationNotes,
			"IsDelivery":        newOrder.IsDelivery,
			"OrderStatus":       newOrder.OrderStatus,
		}

		err = c.Email.SendEmail(
			c.Log,
			[]string{newOrder.Email},
			[]string{},
			mail_subject_cust,
			baseTemplatePath,
			childPath,
			data,
		)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// send admin
		childPath = fmt.Sprintf("../internal/templates/%s/email/order_status_admin.html", request.Lang)
		err = c.Email.SendEmail(
			c.Log,
			[]string{newApp.Email},
			[]string{},
			mail_subject_admin,
			baseTemplatePath,
			childPath,
			data,
		)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
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

func (c *OrderUseCase) GetInvoice(ctx context.Context, invoiceId string) (*model.OrderResponse, *model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx)

	newOrder := new(entity.Order)
	newOrder.Invoice = invoiceId
	if err := c.OrderRepository.FindOrderByInvoiceId(tx, newOrder, invoiceId); err != nil {
		c.Log.Warnf("failed to get order by invoice id : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get order by invoice id : %+v", err))
	}

	newApplication := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApplication); err != nil {
		c.Log.Warnf("failed to get application setting : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get application setting : %+v", err))
	}

	return converter.OrderToResponse(newOrder), converter.ApplicationToResponse(newApplication), nil
}
