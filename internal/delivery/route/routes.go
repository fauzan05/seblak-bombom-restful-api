package route

import (
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/delivery/http"
	xenditController "seblak-bombom-restful-api/internal/delivery/http/xendit"

	"github.com/gofiber/fiber/v2"
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
	MidtransSnapOrderController       *http.MidtransSnapOrderController
	MidtransCoreAPIOrderController    *http.MidtransCoreAPIOrderController
	XenditQRCodeTransactionController *xenditController.XenditQRCodeTransctionController
	XenditCallbackController          *xenditController.XenditCallbackController
	XenditPayoutController            *xenditController.XenditPayoutController
	PayoutController                  *http.PayoutController
	ApplicationController             *http.ApplicationController
	CartController                    *http.CartController
	AuthMiddleware                    fiber.Handler
	RoleMiddleware                    fiber.Handler
	AuthXenditMiddleware              fiber.Handler
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
	
	// User
	api.Post("/users", c.UserController.Register)
	api.Post("/users/login", c.UserController.Login)

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

	// Midtrans
	api.Get("/midtrans/snap/orders/notification", c.MidtransSnapOrderController.GetSnapOrderNotification)
	api.Get("/midtrans/core-api/orders/notification", c.MidtransCoreAPIOrderController.GetCoreAPIOrderNotification)

	// Images
	uploadsDir := "../uploads/images"
	api.Static("/uploads", uploadsDir)

	api.Get("/image/:dir/:filename", func(c *fiber.Ctx) error {
		dir := c.Params("dir")           // Direktori (contoh: products, applications)
		filename := c.Params("filename") // Nama file

		// Gabungkan path menggunakan filepath.Join untuk keamanan
		filepath := filepath.Join(uploadsDir, dir, filename)

		// Mengecek apakah file ada di direktori uploads
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			// Jika file tidak ditemukan, kembalikan error 404
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		// Kirimkan gambar jika ditemukan
		return c.SendFile(filepath)
	})

	api.Get("/applications", c.ApplicationController.Get)
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
	auth.Get("/orders/users/current", c.OrderController.GetAllCurrent)
	auth.Get("/orders/users/:userId", c.OrderController.GetAllByUserId)

	// Product review
	auth.Post("/reviews", c.ProductReviewController.Create)

	// Midtrans
	api.Post("/midtrans/snap/orders", c.MidtransSnapOrderController.CreateSnap)
	api.Post("/midtrans/core-api/orders", c.MidtransCoreAPIOrderController.CreateCoreAPI)
	api.Get("/midtrans/core-api/orders/:orderId", c.MidtransCoreAPIOrderController.GetCoreAPIOrder)

	// Xendit
	api.Post("/xendit/orders/qr-code/transaction", c.XenditQRCodeTransactionController.Create)
	api.Get("/xendit/orders/:orderId/qr-code/transaction", c.XenditQRCodeTransactionController.GetTransaction)
	api.Post("/xendit/payout-request/:payoutId/cancel", c.XenditPayoutController.Cancel)
	api.Get("/xendit/payout-request/:payoutId", c.XenditPayoutController.GetPayoutById)

	// Cart
	api.Post("/carts", c.CartController.Create)
	api.Get("/carts", c.CartController.GetAllCurrent)
	api.Patch("/carts/cart-items/:cartItemId", c.CartController.Update)
	api.Delete("/carts/cart-items/:cartItemId", c.CartController.Delete)

	// Order
	auth.Patch("/orders/:orderId/status", c.OrderController.UpdateOrderStatus)

	// Xendit payout
	auth.Post("/xendit/payouts/:userId", c.XenditPayoutController.Create)

	// Payout
	auth.Post("/payouts/:userId", c.PayoutController.Create)
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
}
