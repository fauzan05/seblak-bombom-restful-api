package config

import (
	"seblak-bombom-restful-api/internal/delivery/http"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/delivery/route"
	"seblak-bombom-restful-api/internal/repository"
	"seblak-bombom-restful-api/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB            *gorm.DB
	App           *fiber.App
	Log           *logrus.Logger
	Validate      *validator.Validate
	Config        *viper.Viper
	SnapClient    *snap.Client
	CoreAPIClient *coreapi.Client
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
	deliveryRepository := repository.NewDeliveryRepository(config.Log)
	productReviewRepository := repository.NewProductReviewRepository(config.Log)
	orderProductRepository := repository.NewOrderProductRepository(config.Log)
	midtransSnapOrderRepository := repository.NewMidtransSnapOrderRepository(config.Log)
	applicationRepository := repository.NewApplicationRepository(config.Log)
	cartRepository := repository.NewCartRepository(config.Log)
	cartItemRepository := repository.NewCartItemRepository(config.Log)

	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenRepository, addressRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, userRepository, addressRepository, userUseCase)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validate, categoryRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, config.Validate, categoryRepository, productRepository, imageRepository)
	imageUseCase := usecase.NewImageUseCase(config.DB, config.Log, config.Validate, imageRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, productRepository, categoryRepository, addressRepository, discountRepository, deliveryRepository, orderProductRepository)
	discountUseCase := usecase.NewDiscountUseCase(config.DB, config.Log, config.Validate, discountRepository)
	deliveryUseCase := usecase.NewDeliveryUseCase(config.DB, config.Log, config.Validate, deliveryRepository)
	productReviewUseCase := usecase.NewProductReviewUseCase(config.DB, config.Log, config.Validate, productReviewRepository)
	midtransSnapOrderUseCase := usecase.NewMidtransSnapOrderUseCase(config.Log, config.Validate, orderRepository, config.SnapClient, config.CoreAPIClient, config.DB, midtransSnapOrderRepository, productRepository)
	applicationUseCase := usecase.NewApplicationUseCase(config.DB, config.Log, config.Validate, applicationRepository)
	cartUseCase := usecase.NewCartUseCase(config.DB, config.Log, config.Validate, cartRepository, productRepository, cartItemRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)
	productController := http.NewProductController(productUseCase, config.Log)
	imageController := http.NewImageController(imageUseCase, config.Log)
	orderController := http.NewOrderController(orderUseCase, config.Log)
	discountController := http.NewDiscountController(discountUseCase, config.Log)
	deliveryController := http.NewDeliveryController(deliveryUseCase, config.Log)
	productReviewController := http.NewProductReviewController(productReviewUseCase, config.Log)
	midtransSnapOrderController := http.NewMidtransController(midtransSnapOrderUseCase, config.Log)
	applicationController := http.NewApplicationController(applicationUseCase, config.Log)
	cartController := http.NewCartController(cartUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	roleMiddleware := middleware.NewRole(userUseCase)

	routeConfig := route.RouteConfig{
		App:                         config.App,
		UserController:              userController,
		AddressController:           addressController,
		CategoryController:          categoryController,
		ProductController:           productController,
		ImageController:             imageController,
		OrderController:             orderController,
		DiscountController:          discountController,
		DeliveryController:          deliveryController,
		ProductReviewController:     productReviewController,
		MidtransSnapOrderController: midtransSnapOrderController,
		ApplicationController:       applicationController,
		CartController:              cartController,
		AuthMiddleware:              authMiddleware,
		RoleMiddleware:              roleMiddleware,
	}
	routeConfig.Setup()
}
