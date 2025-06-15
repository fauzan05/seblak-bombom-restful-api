package route

import (
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/delivery/http"
	xenditController "seblak-bombom-restful-api/internal/delivery/http/xendit"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pusher/pusher-http-go/v5"
)

type RouteConfig struct {
	App                               *fiber.App
	UserController                    *http.UserController
	AddressController                 *http.AddressController
	CategoryController                *http.CategoryController
	ProductController                 *http.ProductController
	OrderController                   *http.OrderController
	DiscountCouponController          *http.DiscountCouponController
	DeliveryController                *http.DeliveryController
	ProductReviewController           *http.ProductReviewController
	XenditQRCodeTransactionController *xenditController.XenditQRCodeTransctionController
	XenditCallbackController          *xenditController.XenditCallbackController
	XenditPayoutController            *xenditController.XenditPayoutController
	PayoutController                  *http.PayoutController
	ApplicationController             *http.ApplicationController
	CartController                    *http.CartController
	WalletController                  *http.WalletController
	AuthMiddleware                    fiber.Handler
	RoleMiddleware                    fiber.Handler
	AuthXenditMiddleware              fiber.Handler
	AuthAdminCreationMiddleware       fiber.Handler
	PusherClient                      pusher.Client
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupXenditCallbacksRoute()
	c.SetupAuthRoute()
	c.SetupAuthAdminRoute()
}

func (c *RouteConfig) SetupXenditCallbacksRoute() {
	api := c.App.Group("/api")
	authToken := api.Use(c.AuthXenditMiddleware)
	// Xendit QR Code Callback
	authToken.Post("/xendits/payment-request/notifications/callback", c.XenditCallbackController.GetPaymentRequestCallbacks)
	authToken.Post("/xendits/payout-request/notifications/callback", c.XenditCallbackController.GetPayoutRequestCallbacks)
}

// GUEST
func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("/api")
	wd, _ := os.Getwd()
	staticPath := filepath.Join(wd, "../internal/templates/assets")
	api.Static("/assets", staticPath)

	// User
	api.Post("/users/register", c.UserController.Register)
	api.Post("/users/login", c.UserController.Login)
	api.Post("/users/forgot-password", c.UserController.CreateForgotPassword)
	api.Post("/users/forgot-password/:passwordResetId/validate", c.UserController.ValidateForgotPassword)
	api.Post("/users/forgot-password/:passwordResetId/reset-password", c.UserController.ResetPassword)
	api.Get("/users/verify-email/:token", c.UserController.VerifyEmailRegistration)
	c.App.Get("/verified-success/:token", c.UserController.ShowVerifiedSuccess)
	c.App.Get("/verified-failed/:token", c.UserController.ShowVerifiedFailed)

	// Discount Coupon
	api.Get("/discount-coupons", c.DiscountCouponController.GetAll)
	api.Get("/discount-coupons/:discountId", c.DiscountCouponController.Get)

	// Category
	api.Get("/categories/:categoryId", c.CategoryController.Get)
	api.Get("/categories", c.CategoryController.GetAll)

	// Delivery
	api.Get("/deliveries", c.DeliveryController.GetAll)

	// Product
	api.Get("/products", c.ProductController.GetAll)
	api.Get("/products/:productId", c.ProductController.Get)

	// Images
	uploadsDir := "../uploads/images"
	api.Get("/image/*", func(c *fiber.Ctx) error {
		relativePath := c.Params("*") // Menangkap seluruh path setelah /image/
		if relativePath == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing file path")
		}

		// Gabungkan path dengan root upload
		unsafePath := filepath.Join(uploadsDir, relativePath)
		cleanPath := filepath.Clean(unsafePath)

		// Pastikan cleanPath masih dalam uploadsDir
		absUploadsDir, _ := filepath.Abs(uploadsDir)
		absCleanPath, _ := filepath.Abs(cleanPath)
		if !strings.HasPrefix(absCleanPath, absUploadsDir) {
			return c.Status(fiber.StatusForbidden).SendString("Access denied")
		}

		// Cek apakah file ada
		if _, err := os.Stat(absCleanPath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		// Kirim file
		return c.SendFile(absCleanPath)
	})

	// Application
	api.Get("/applications", c.ApplicationController.Get)
	api.Use(c.AuthAdminCreationMiddleware).Post("/applications-use-admin-key", c.ApplicationController.Create) // add & update

	api.Get("/test-pusher", func(f *fiber.Ctx) error {
		message := f.Query("message", "")
		data := map[string]string{"message": message}
		err := c.PusherClient.Trigger("seblak_bombom_api_channel", "event_testing1", data)
		if err != nil {
			return err
		}
		return f.Status(fiber.StatusOK).JSON(data)
	})
}

// USER
func (c *RouteConfig) SetupAuthRoute() {
	api := c.App.Group("/api")
	auth := api.Use(c.AuthMiddleware)

	// User
	auth.Get("/users/current", c.UserController.GetCurrent)
	auth.Patch("/users/current", c.UserController.Update)
	auth.Delete("/users/logout", c.UserController.Logout)
	auth.Delete("/users/current", c.UserController.RemoveAccount)
	auth.Patch("/users/current/password", c.UserController.UpdatePassword)

	// Address
	auth.Post("/users/current/addresses", c.AddressController.Add)
	auth.Get("/users/current/addresses", c.AddressController.GetAll)
	auth.Get("/users/current/addresses/:addressId", c.AddressController.Get)
	auth.Put("/users/current/addresses/:addressId", c.AddressController.Update)
	auth.Delete("/users/current/addresses", c.AddressController.Remove)

	// Order
	auth.Post("/orders", c.OrderController.Create)
	auth.Get("/orders/:orderId", c.OrderController.GetOrderById)
	auth.Get("/orders/users/:userId", c.OrderController.GetAllByUserId)
	auth.Patch("/orders/:orderId/status", c.OrderController.UpdateOrderStatus)
	auth.Get("/orders", c.OrderController.GetAll)
	auth.Get("/orders/:invoiceId/invoice", c.OrderController.ShowInvoiceByOrderId)

	// Product review
	auth.Post("/reviews", c.ProductReviewController.Create)

	// Xendit
	auth.Post("/xendit/orders/qr-code/transaction", c.XenditQRCodeTransactionController.Create)
	auth.Get("/xendit/orders/:orderId/qr-code/transaction", c.XenditQRCodeTransactionController.GetTransaction)
	auth.Post("/xendit/payout-request/:payoutId/cancel", c.XenditPayoutController.Cancel)
	auth.Get("/xendit/payout-request/:payoutId", c.XenditPayoutController.GetPayoutById)

	// Cart
	auth.Post("/carts", c.CartController.Create)
	auth.Get("/carts", c.CartController.GetAllCurrent)
	auth.Patch("/carts/cart-items/:cartItemId", c.CartController.Update)
	auth.Delete("/carts/cart-items/:cartItemId", c.CartController.Delete)

	// Xendit payout
	auth.Post("/xendit/payouts/:userId", c.XenditPayoutController.Create)

	// Payout
	auth.Post("/payouts/:userId", c.PayoutController.Create)

	// Wallet
	auth.Post("/wallets/withdraw-cust", c.WalletController.WithdrawCustRequest)
}

// ADMIN
func (c *RouteConfig) SetupAuthAdminRoute() {
	api := c.App.Group("/api")
	auth := api.Use(c.AuthMiddleware, c.RoleMiddleware)

	// Category
	auth.Post("/categories", c.CategoryController.Create)
	auth.Put("/categories/:categoryId", c.CategoryController.Edit)
	auth.Delete("/categories", c.CategoryController.Remove)

	// Product
	auth.Post("/products", c.ProductController.Create)
	auth.Put("/products/:productId", c.ProductController.Edit)
	auth.Delete("/products", c.ProductController.Remove)

	// Discount
	auth.Post("/discount-coupons", c.DiscountCouponController.Create)
	auth.Put("/discount-coupons/:discountId", c.DiscountCouponController.Update)
	auth.Delete("/discount-coupons", c.DiscountCouponController.Delete)

	// Delivery
	auth.Post("/deliveries", c.DeliveryController.Create)
	auth.Put("/deliveries/:deliveryId", c.DeliveryController.Update)
	auth.Delete("/deliveries", c.DeliveryController.Remove)

	// Application
	auth.Post("/applications", c.ApplicationController.Create) // add & update

	// Balance
	auth.Get("/balance", c.XenditPayoutController.GetAdminBalance)

	// Wallet
	auth.Patch("/wallets/:withdrawRequestId/withdraw-approval", c.WalletController.WithdrawAdminApproval)
}
