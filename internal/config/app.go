package config

import (
	"seblak-bombom-restful-api/internal/delivery/http"
	xenditController "seblak-bombom-restful-api/internal/delivery/http/xendit"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/delivery/route"
	"seblak-bombom-restful-api/internal/repository"
	"seblak-bombom-restful-api/internal/usecase"
	xenditUseCase "seblak-bombom-restful-api/internal/usecase/xendit"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xendit/xendit-go/v6"
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
	XenditClient  *xendit.APIClient
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
	discountCouponRepository := repository.NewDiscountRepository(config.Log)
	deliveryRepository := repository.NewDeliveryRepository(config.Log)
	productReviewRepository := repository.NewProductReviewRepository(config.Log)
	orderProductRepository := repository.NewOrderProductRepository(config.Log)
	midtransSnapOrderRepository := repository.NewMidtransSnapOrderRepository(config.Log)
	midtransCoreAPIOrderRepository := repository.NewMidtransCoreAPIOrderRepository(config.Log)
	xenditTransactionRepository := repository.NewXenditTransactionRepository(config.Log)
	applicationRepository := repository.NewApplicationRepository(config.Log)
	cartRepository := repository.NewCartRepository(config.Log)
	cartItemRepository := repository.NewCartItemRepository(config.Log)
	walletRepository := repository.NewWalletRepository(config.Log)

	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenRepository, addressRepository, walletRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, userRepository, addressRepository, deliveryRepository, userUseCase)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validate, categoryRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, config.Validate, categoryRepository, productRepository, imageRepository)
	imageUseCase := usecase.NewImageUseCase(config.DB, config.Log, config.Validate, imageRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, productRepository, categoryRepository, addressRepository, discountCouponRepository, deliveryRepository, orderProductRepository, walletRepository)
	discountCouponUseCase := usecase.NewDiscountCouponUseCase(config.DB, config.Log, config.Validate, discountCouponRepository)
	deliveryUseCase := usecase.NewDeliveryUseCase(config.DB, config.Log, config.Validate, deliveryRepository)
	productReviewUseCase := usecase.NewProductReviewUseCase(config.DB, config.Log, config.Validate, productReviewRepository)
	midtransSnapOrderUseCase := usecase.NewMidtransSnapOrderUseCase(config.Log, config.Validate, orderRepository, config.SnapClient, config.CoreAPIClient, config.DB, midtransSnapOrderRepository, productRepository)
	midtransCoreApiOrderUseCase := usecase.NewMidtransCoreAPIOrderUseCase(config.Log, config.Validate, orderRepository, config.CoreAPIClient, config.DB, midtransCoreAPIOrderRepository, productRepository)
	xenditTransactionQRCodeUseCase := xenditUseCase.NewXenditTransactionQRCodeUseCase(config.DB, config.Log, config.Validate, orderRepository, xenditTransactionRepository, config.XenditClient)
	applicationUseCase := usecase.NewApplicationUseCase(config.DB, config.Log, config.Validate, applicationRepository)
	cartUseCase := usecase.NewCartUseCase(config.DB, config.Log, config.Validate, cartRepository, productRepository, cartItemRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)
	productController := http.NewProductController(productUseCase, config.Log)
	imageController := http.NewImageController(imageUseCase, config.Log)
	orderController := http.NewOrderController(orderUseCase, config.Log)
	discountCouponController := http.NewDiscountCouponController(discountCouponUseCase, config.Log)
	deliveryController := http.NewDeliveryController(deliveryUseCase, config.Log)
	productReviewController := http.NewProductReviewController(productReviewUseCase, config.Log)
	midtransSnapOrderController := http.NewMidtransSnapOrderController(midtransSnapOrderUseCase, config.Log)
	midtransCoreAPIOrderController := http.NewMidtransCoreAPIOrderController(midtransCoreApiOrderUseCase, config.Log)
	xenditQRCodeTransactionController := xenditController.NewXenditQRCodeTransctionController(xenditTransactionQRCodeUseCase, config.Log)
	applicationController := http.NewApplicationController(applicationUseCase, config.Log)
	cartController := http.NewCartController(cartUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	roleMiddleware := middleware.NewRole(userUseCase)

	routeConfig := route.RouteConfig{
		App:                               config.App,
		UserController:                    userController,
		AddressController:                 addressController,
		CategoryController:                categoryController,
		ProductController:                 productController,
		ImageController:                   imageController,
		OrderController:                   orderController,
		DiscountCouponController:          discountCouponController,
		DeliveryController:                deliveryController,
		ProductReviewController:           productReviewController,
		MidtransSnapOrderController:       midtransSnapOrderController,
		MidtransCoreAPIOrderController:    midtransCoreAPIOrderController,
		XenditQRCodeTransactionController: xenditQRCodeTransactionController,
		ApplicationController:             applicationController,
		CartController:                    cartController,
		AuthMiddleware:                    authMiddleware,
		RoleMiddleware:                    roleMiddleware,
	}
	routeConfig.Setup()
}
