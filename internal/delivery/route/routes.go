package route

import (
	"seblak-bombom-restful-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	UserController     *http.UserController
	AddressController  *http.AddressController
	CategoryController *http.CategoryController
	ProductController  *http.ProductController
	AuthMiddleware     fiber.Handler
	RoleMiddleware     fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
	c.SetupAuthAdminRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/login", c.UserController.Login)
}

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

	// Category
	auth.Get("/categories/:categoryId", c.CategoryController.Get)
	auth.Get("/categories", c.CategoryController.GetAll)

	// Product
	auth.Get("/products", c.ProductController.GetAll)
	auth.Get("/products/:productId", c.ProductController.Get)
}

func (c *RouteConfig) SetupAuthAdminRoute() {
	api := c.App.Group("/api")
	auth := api.Use(c.AuthMiddleware, c.RoleMiddleware)

	// category
	auth.Post("/categories", c.CategoryController.Create)
	auth.Put("/categories/:categoryId", c.CategoryController.Edit)
	auth.Delete("/categories/:categoryId", c.CategoryController.Remove)

	// Product
	auth.Post("/products", c.ProductController.Create)
	auth.Put("/products/:productId", c.ProductController.Edit)
	auth.Delete("/products/:productId", c.ProductController.Remove)
}
