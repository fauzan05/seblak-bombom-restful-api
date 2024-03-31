package route

import (
	"seblak-bombom-restful-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	AddressController *http.AddressController
	CategoryController *http.CategoryController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	// User
	c.App.Get("/api/users/current", c.UserController.GetCurrent)
	c.App.Patch("/api/users/current", c.UserController.Update)
	c.App.Delete("/api/users/logout", c.UserController.Logout)
	c.App.Delete("/api/users/current", c.UserController.RemoveAccount)
	c.App.Patch("/api/users/current/password", c.UserController.UpdatePassword)
	
	// Address
	c.App.Post("/api/users/current/addresses", c.AddressController.Add)
	c.App.Get("/api/users/current/addresses", c.AddressController.GetAll)
	c.App.Get("/api/addresses/:addressId", c.AddressController.Get)
	c.App.Put("/api/addresses/:addressId", c.AddressController.Update)
	c.App.Delete("/api/addresses/:addressId", c.AddressController.Remove)

	// Category
	c.App.Post("/api/categories", c.CategoryController.Create)
	c.App.Get("/api/categories/:categoryId", c.CategoryController.Get)
	c.App.Get("/api/categories", c.CategoryController.GetAll)
	c.App.Put("/api/categories/:categoryId", c.CategoryController.Edit)
	c.App.Delete("/api/categories/:categoryId", c.CategoryController.Remove)
	

}
