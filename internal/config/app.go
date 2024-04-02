package config

import (
	"seblak-bombom-restful-api/internal/delivery/http"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/delivery/route"
	"seblak-bombom-restful-api/internal/repository"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	tokenRepository := repository.NewTokenRepository(config.Log)
	addressRepository := repository.NewAddressRepository(config.Log)
	categoryRepository := repository.NewCategoryRepository(config.Log)
	productRepository := repository.NewProductRepository(config.Log)
	imageRepository := repository.NewImageRepository(config.Log)
	orderRepository := repository.NewOrderRepository(config.Log)
	discountRepository := repository.NewDiscountRepository(config.Log)

	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenRepository, addressRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, userRepository, addressRepository, userUseCase)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validate, categoryRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, config.Validate, categoryRepository, productRepository)
	imageUseCase := usecase.NewImageUseCase(config.DB, config.Log, config.Validate, imageRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, productRepository, categoryRepository, addressRepository)
	discountUseCase := usecase.NewDiscountUseCase(config.DB, config.Log, config.Validate, discountRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)
	productController := http.NewProductController(productUseCase, config.Log)
	imageController := http.NewImageController(imageUseCase, config.Log)
	orderController := http.NewOrderController(orderUseCase, config.Log)
	discountController := http.NewDiscountController(discountUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	roleMiddleware := middleware.NewRole(userUseCase)

	routeConfig := route.RouteConfig{
		App:                config.App,
		UserController:     userController,
		AddressController:  addressController,
		CategoryController: categoryController,
		ProductController:  productController,
		ImageController:    imageController,
		OrderController:    orderController,
		DiscountController: discountController,
		AuthMiddleware:     authMiddleware,
		RoleMiddleware:     roleMiddleware,
	}
	routeConfig.Setup()
}
