package route

import (
	"fmt"
	"os"
	"seblak-bombom-restful-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                         *fiber.App
	UserController              *http.UserController
	AddressController           *http.AddressController
	CategoryController          *http.CategoryController
	ProductController           *http.ProductController
	ImageController             *http.ImageController
	OrderController             *http.OrderController
	DiscountController          *http.DiscountController
	DeliveryController          *http.DeliveryController
	ProductReviewController     *http.ProductReviewController
	MidtransSnapOrderController *http.MidtransSnapOrderController
	ApplicationController       *http.ApplicationController
	CartController              *http.CartController
	AuthMiddleware              fiber.Handler
	RoleMiddleware              fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
	c.SetupAuthAdminRoute()
}

// GUEST
func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("/api")
	api.Post("/users", c.UserController.Register)
	api.Post("/users/login", c.UserController.Login)
	api.Get("/discounts", c.DiscountController.GetAll)
	api.Get("/discounts/:discountId", c.DiscountController.Get)

	// Category
	api.Get("/categories/:categoryId", c.CategoryController.Get)
	api.Get("/categories", c.CategoryController.GetAll)

	// delivery
	api.Get("/deliveries", c.DeliveryController.Get)

	// Product
	api.Get("/products", c.ProductController.GetAll)
	api.Get("/products/:productId", c.ProductController.Get)

	// Midtrans
	api.Get("/midtrans/snap/orders/notification", c.MidtransSnapOrderController.GetSnapOrderNotification)

	// Images
	uploadsDir := "../uploads/images/products/"
	api.Static("/uploads", uploadsDir)
	api.Get("/image/:filename", func (c *fiber.Ctx) error  {
		filename := c.Params("filename")
		filepath := fmt.Sprintf("%s/%s", uploadsDir, filename)

		// Mengecek apakah file ada di direktori uploads
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			// Jika file tidak ditemukan, kembalikan error 404
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		// Kirimkan gambar jika ditemukan
		return c.SendFile(filepath)
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
	auth.Get("/addresses/:addressId", c.AddressController.Get)
	auth.Put("/addresses/:addressId", c.AddressController.Update)
	auth.Delete("/addresses/:addressId", c.AddressController.Remove)

	// order
	auth.Post("/orders", c.OrderController.Create)
	auth.Get("/orders/users/current", c.OrderController.GetAllCurrent)
	auth.Get("/orders/users/:userId", c.OrderController.GetAllByUserId)

	// product review
	auth.Post("/reviews", c.ProductReviewController.Create)

	// Midtrans
	api.Post("/midtrans/snap/orders", c.MidtransSnapOrderController.CreateSnap)

	// Cart
	api.Post("/carts", c.CartController.Create)
	api.Get("/carts", c.CartController.GetAllCurrent)

	// order
	auth.Patch("/orders/:orderId", c.OrderController.Update)
}

// ADMIN
func (c *RouteConfig) SetupAuthAdminRoute() {
	api := c.App.Group("/api")
	auth := api.Use(c.AuthMiddleware, c.RoleMiddleware)

	// category
	auth.Post("/categories", c.CategoryController.Create)
	auth.Put("/categories/:categoryId", c.CategoryController.Edit)
	auth.Delete("/categories", c.CategoryController.Remove)

	// Product
	auth.Post("/products", c.ProductController.Create)
	auth.Put("/products/:productId", c.ProductController.Edit)
	auth.Delete("/products", c.ProductController.Remove)

	// image
	auth.Post("/images", c.ImageController.Creates)
	auth.Put("/images", c.ImageController.EditPosition)
	auth.Delete("/images/:imageId", c.ImageController.Remove)

	// discount
	auth.Post("/discounts", c.DiscountController.Create)
	auth.Put("/discounts/:discountId", c.DiscountController.Update)
	auth.Delete("/discounts/:discountId", c.DiscountController.Delete)

	// delivery
	auth.Post("/deliveries", c.DeliveryController.Create)
	auth.Put("/deliveries/:deliveryId", c.DeliveryController.Update)
	auth.Delete("/deliveries/:deliveryId", c.DeliveryController.Remove)

	// application
	auth.Post("/applications", c.ApplicationController.Create)
	auth.Put("/applications", c.ApplicationController.Update)
	auth.Get("/applications", c.ApplicationController.Get)
}
