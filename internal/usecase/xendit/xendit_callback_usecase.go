package usecase

import (
	"fmt"
	"html/template"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payment_request"
	"gorm.io/gorm"
)

type XenditCallbackUseCase struct {
	DB                          *gorm.DB
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	XenditClient                *xendit.APIClient
	OrderRepository             *repository.OrderRepository
	XenditTransactionRepository *repository.XenditTransctionRepository
	UserRepository              *repository.UserRepository
	WalletRepository            *repository.WalletRepository
	XenditPayoutRepository      *repository.XenditPayoutRepository
	PayoutRepository            *repository.PayoutRepository
	ApplicationRepository       *repository.ApplicationRepository
	NotificationRepository      *repository.NotificationRepository
	Email                       *mailer.EmailWorker
}

func NewXenditCallbackUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, xenditTransactionRepository *repository.XenditTransctionRepository,
	xenditClient *xendit.APIClient, xenditPayoutRepository *repository.XenditPayoutRepository,
	userRepository *repository.UserRepository, walletRepository *repository.WalletRepository,
	payoutRepository *repository.PayoutRepository, applicationRepository *repository.ApplicationRepository,
	notificationRepository *repository.NotificationRepository, email *mailer.EmailWorker) *XenditCallbackUseCase {
	return &XenditCallbackUseCase{
		DB:                          db,
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		XenditTransactionRepository: xenditTransactionRepository,
		XenditPayoutRepository:      xenditPayoutRepository,
		XenditClient:                xenditClient,
		UserRepository:              userRepository,
		WalletRepository:            walletRepository,
		PayoutRepository:            payoutRepository,
		ApplicationRepository:       applicationRepository,
		NotificationRepository:      notificationRepository,
		Email:                       email,
	}
}

func (c *XenditCallbackUseCase) UpdateStatusPaymentRequestCallback(ctx *fiber.Ctx, request *model.XenditGetPaymentRequestCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newXenditTransaction := new(entity.XenditTransactions)
	count, err := c.XenditTransactionRepository.FindXenditTransaction(tx, newXenditTransaction, request.Data.PaymentMethod.ID)
	if err != nil {
		c.Log.Warnf("failed to get xendit transaction from database : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get xendit transaction from database : %+v", err))
	}

	if count > 0 {
		// update datanya
		if newXenditTransaction.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt
			status := request.Data.Status
			orderId := newXenditTransaction.OrderId
			updateXenditTransaction := map[string]any{
				"status":     status,
				"updated_at": updatedAt.ToTime(),
			}

			*newXenditTransaction = entity.XenditTransactions{
				ID: newXenditTransaction.ID,
			}

			if err := c.XenditTransactionRepository.UpdateCustomColumns(tx, newXenditTransaction, updateXenditTransaction); err != nil {
				c.Log.Warnf("failed to update xendit transaction status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update xendit transaction status into database : %+v", err))
			}

			var payment_status helper.PaymentStatus
			var is_send_email bool
			var email_subject string
			if status == string(payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED) {
				is_send_email = true
				payment_status = helper.PAID_PAYMENT
				email_subject = "Payment Successfull"
				if request.Lang == helper.INDONESIA {
					email_subject = "Pembayaran Berhasil"
				}
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_CANCELED) {
				is_send_email = true
				payment_status = helper.CANCELLED_PAYMENT
				email_subject = "Payment Cancelled"
				if request.Lang == helper.INDONESIA {
					email_subject = "Pembayaran Dibatalkan"
				}
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_FAILED) {
				is_send_email = true
				payment_status = helper.FAILED_PAYMENT
				email_subject = "Payment Failed"
				if request.Lang == helper.INDONESIA {
					email_subject = "Pembayaran Gagal"
				}
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_EXPIRED) {
				is_send_email = true
				payment_status = helper.EXPIRED_PAYMENT
				email_subject = "Payment Expired"
				if request.Lang == helper.INDONESIA {
					email_subject = "Pembayaran Kadaluwarsa"
				}
			}

			if status == string(payment_request.PAYMENTREQUESTSTATUS_PENDING) {
				payment_status = helper.PENDING_PAYMENT
			}

			updateOrderStatus := map[string]any{
				"payment_status": payment_status,
				"updated_at":     updatedAt.ToTime(),
			}

			newOrder := new(entity.Order)
			newOrder.ID = orderId
			if err := c.OrderRepository.UpdateCustomColumns(tx, newOrder, updateOrderStatus); err != nil {
				c.Log.Warnf("failed to update order status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update order status into database : %+v", err))
			}

			if is_send_email {
				if err := c.OrderRepository.FindWith2Preloads(tx, newOrder, "OrderProducts", "OrderProducts.Product"); err != nil {
					c.Log.Warnf("failed to find newly created order : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find newly created order : %+v", err))
				}

				productsSelected := []map[string]any{}
				for _, product := range newOrder.OrderProducts {
					var productImageBase64 string
					productImagePath := fmt.Sprintf("../uploads/images/products/%s", product.ProductFirstImagePosition)
					productImageBase64, err := helper.ImageToBase64(productImagePath)
					if err != nil {
						c.Log.Warnf("failed to convert product image to base64 : %+v", err)
						return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert product image to base64 : %+v", err))
					}

					productImage := map[string]any{
						"ProductImageFilename": product.ProductFirstImagePosition,
						"ProductImage":         productImageBase64,
						"ProductName":          product.ProductName,
						"Quantity":             product.Quantity,
						"Price":                helper.FormatNumberFloat32(product.Price),
					}

					productsSelected = append(productsSelected, productImage)
				}

				newApp := new(entity.Application)
				if err := c.ApplicationRepository.FindFirst(tx, newApp); err != nil {
					c.Log.Warnf("failed to find application from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application from database : %+v", err))
				}
				if newApp.LogoFilename == "" {
					c.Log.Warnf("application logo has not uploaded yet!")
					return fiber.NewError(fiber.StatusBadRequest, "application logo has not uploaded yet!")
				}

				logoImagePath := fmt.Sprintf("../uploads/images/application/%s", newApp.LogoFilename)
				logoImageBase64, err := helper.ImageToBase64(logoImagePath)
				if err != nil {
					c.Log.Warnf("failed to convert logo to base64 : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to convert logo to base64 : %+v", err))
				}

				newMail := new(model.Mail)
				newMail.To = []string{newOrder.Email}
				newMail.Subject = email_subject
				baseTemplatePath := "../internal/templates/base_template_email1.html"
				childPath := fmt.Sprintf("../internal/templates/%s/email/order_payment.html", request.Lang)
				tmpl, err := template.ParseFiles(baseTemplatePath, childPath)
				if err != nil {
					c.Log.Warnf("failed to parse template file html : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
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
					"TotalAmount":      helper.FormatNumberFloat32(newOrder.TotalFinalPrice),
					"Year":             time.Now().Format("2006"),
					"CustomerNotes":    newOrder.Note,
					"ShippingMethod":   newOrder.IsDelivery,
					"ShippingCost":     helper.FormatNumberFloat32(newOrder.DeliveryCost),
					"ServiceFee":       helper.FormatNumberFloat32(newOrder.ServiceFee),
					"Discount":         helper.FormatNumberFloat32(newOrder.TotalDiscount),
					"Subject":          newMail.Subject,
					"PaymentStatus":    newOrder.PaymentStatus,
					"PaymentLink":      paymentLink,
					"OrderTrackingURL": orderTrackingURL,
				})
				if err != nil {
					c.Log.Warnf("failed to execute template file html : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
				}
				newMail.Template = *bodyBuilder
				c.Email.Mailer.SenderName = fmt.Sprintf("System %s", newApp.AppName)
				// send email
				select {
				case c.Email.MailQueue <- *newMail:
				default:
					c.Log.Warnf("email queue full, failed to send to %s", newOrder.Email)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("email queue full, failed to send to %s", newOrder.Email))
				}

				newNotification := new(entity.Notification)
				newNotification.UserID = newOrder.UserId
				newNotification.Title = newMail.Subject
				newNotification.IsRead = false
				newNotification.Type = helper.TRANSACTION
				baseTemplatePath = "../internal/templates/base_template_notification1.html"
				childPath = fmt.Sprintf("../internal/templates/%s/notification/order_payment.html", request.Lang)
				tmpl, err = template.ParseFiles(baseTemplatePath, childPath)
				if err != nil {
					c.Log.Warnf("failed to parse template file html : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to parse template file html : %+v", err))
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
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to execute template file html : %+v", err))
				}

				newNotification.BodyContent = bodyBuilder.String()
				if err := c.NotificationRepository.Create(tx, newNotification); err != nil {
					c.Log.Warnf("failed to create notification into database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create notification into database : %+v", err))
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return nil
}

func (c *XenditCallbackUseCase) UpdateStatusPayoutRequestCallback(ctx *fiber.Ctx, request *model.XenditGetPayoutRequestCallbackStatus) error {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newXenditPayout := new(entity.XenditPayout)
	newXenditPayout.ID = request.Data.PayoutId
	count, err := c.XenditPayoutRepository.FindFirstAndCount(tx, newXenditPayout)
	if err != nil {
		c.Log.Warnf("failed to get xendit transaction from database : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get xendit transaction from database : %+v", err))
	}

	if count > 0 {
		// update datanya
		if newXenditPayout.Status != request.Data.Status {
			// update statusnya
			updatedAt := request.Data.UpdatedAt
			status := request.Data.Status
			updateXenditPayout := map[string]any{
				"status":     status,
				"updated_at": updatedAt,
			}

			*newXenditPayout = entity.XenditPayout{
				ID: newXenditPayout.ID,
			}

			if err := c.XenditPayoutRepository.UpdateCustomColumns(tx, newXenditPayout, updateXenditPayout); err != nil {
				c.Log.Warnf("failed to update xendit payout status into database : %+v", err)
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update xendit payout status into database : %+v", err))
			}

			if request.Data.Status == "SUCCEEDED" {
				// update tb_payout
				newPayout := new(entity.Payout)
				if err := c.PayoutRepository.FindFirstPayoutByXenditPayoutId(tx, newPayout, newXenditPayout.ID); err != nil {
					c.Log.Warnf("failed to get payout by xendit payout id from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get payout by xendit payout id from database : %+v", err))
				}

				if newPayout.ID < 1 {
					c.Log.Warnf("payout not found!")
					return fiber.NewError(fiber.StatusNotFound, "payout not found!")
				} else {
					// update payout
					updateStatus := map[string]any{
						"status": helper.PAYOUT_SUCCEEDED,
					}

					if err := c.PayoutRepository.UpdateCustomColumns(tx, newPayout, updateStatus); err != nil {
						c.Log.Warnf("failed to update payout status in the database : %+v", err)
						return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update payout status in the database : %+v", err))
					}
				}
			}

			if request.Data.Status == "CANCELLED" || request.Data.Status == "failed" || request.Data.Status == "EXPIRED" || request.Data.Status == "REFUNDED" {
				// kembalikan saldonya
				newUser := new(entity.User)
				newUser.ID = newXenditPayout.UserID
				if err := c.UserRepository.FindWithPreloads(tx, newUser, "Wallet"); err != nil {
					c.Log.Warnf("failed to find user wallet from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ailed to find user wallet from database : %+v", err))
				}

				resultBalance := newUser.Wallet.Balance + request.Data.Amount
				// update saldo
				updateBalance := map[string]any{
					"balance": resultBalance,
				}

				newWallet := new(entity.Wallet)
				newWallet.ID = newUser.Wallet.ID
				if err := c.WalletRepository.UpdateCustomColumns(tx, newWallet, updateBalance); err != nil {
					c.Log.Warnf("failed to update wallet balance in the database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance in the database : %+v", err))
				}

				// update tb_payout
				newPayout := new(entity.Payout)
				if err := c.PayoutRepository.FindFirstPayoutByXenditPayoutId(tx, newPayout, newXenditPayout.ID); err != nil {
					c.Log.Warnf("failed to get payout by xendit payout id from database : %+v", err)
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get payout by xendit payout id from database : %+v", err))
				}

				if newPayout.ID < 1 {
					c.Log.Warnf("payout not found!")
					return fiber.NewError(fiber.StatusNotFound, "payout not found!")
				} else {
					status := helper.PAYOUT_CANCELLED
					if request.Data.Status == "failed" {
						status = helper.PAYOUT_FAILED
					} else if request.Data.Status == "EXPIRED" {
						status = helper.PAYOUT_EXPIRED
					} else if request.Data.Status == "REFUNDED" {
						status = helper.PAYOUT_REFUNDED
					}

					// update payout
					updateStatus := map[string]any{
						"status": status,
					}

					if err := c.PayoutRepository.UpdateCustomColumns(tx, newPayout, updateStatus); err != nil {
						c.Log.Warnf("failed to update payout status in the database : %+v", err)
						return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update payout status in the database : %+v", err))
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return nil
}
