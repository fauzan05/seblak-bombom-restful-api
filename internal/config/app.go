package config

import (
	"seblak-bombom-restful-api/internal/delivery/http"
	xenditController "seblak-bombom-restful-api/internal/delivery/http/xendit"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/delivery/route"
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/repository"
	"seblak-bombom-restful-api/internal/usecase"
	xenditUseCase "seblak-bombom-restful-api/internal/usecase/xendit"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
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
	Email         *mailer.EmailWorker
	PDF           *wkhtmltopdf.PDFGenerator
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
	discountUsageRepository := repository.NewDiscountUsageRepository(config.Log)
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
	xenditPayoutRepository := repository.NewXenditPayoutRepository(config.Log)
	payoutRepository := repository.NewPayoutRepository(config.Log)
	notificationRepository := repository.NewNotificationRepository(config.Log)
	passwordResetRepository := repository.NewPasswordResetRepository(config.Log)

	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenRepository, addressRepository, walletRepository, cartRepository, notificationRepository, config.Email, applicationRepository, passwordResetRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, userRepository, addressRepository, deliveryRepository, userUseCase)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validate, categoryRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, config.Validate, categoryRepository, productRepository, imageRepository)
	discountCouponUseCase := usecase.NewDiscountCouponUseCase(config.DB, config.Log, config.Validate, discountCouponRepository)
	deliveryUseCase := usecase.NewDeliveryUseCase(config.DB, config.Log, config.Validate, deliveryRepository)
	productReviewUseCase := usecase.NewProductReviewUseCase(config.DB, config.Log, config.Validate, productReviewRepository)
	midtransSnapOrderUseCase := usecase.NewMidtransSnapOrderUseCase(config.Log, config.Validate, orderRepository, config.SnapClient, config.CoreAPIClient, config.DB, midtransSnapOrderRepository, productRepository)
	midtransCoreApiOrderUseCase := usecase.NewMidtransCoreAPIOrderUseCase(config.Log, config.Validate, orderRepository, config.CoreAPIClient, config.DB, midtransCoreAPIOrderRepository, productRepository)
	xenditTransactionQRCodeUseCase := xenditUseCase.NewXenditTransactionQRCodeUseCase(config.DB, config.Log, config.Validate, orderRepository, xenditTransactionRepository, config.XenditClient)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, productRepository, categoryRepository, addressRepository, discountCouponRepository, discountUsageRepository, deliveryRepository, orderProductRepository, walletRepository, xenditTransactionRepository, xenditTransactionQRCodeUseCase, config.XenditClient, applicationRepository, config.Email)
	applicationUseCase := usecase.NewApplicationUseCase(config.DB, config.Log, config.Validate, applicationRepository)
	cartUseCase := usecase.NewCartUseCase(config.DB, config.Log, config.Validate, cartRepository, productRepository, cartItemRepository)
	xenditCallbackUseCase := xenditUseCase.NewXenditCallbackUseCase(config.DB, config.Log, config.Validate, orderRepository, xenditTransactionRepository, config.XenditClient, xenditPayoutRepository, userRepository, walletRepository, payoutRepository)
	xenditPayoutUseCase := xenditUseCase.NewXenditPayoutUseCase(config.DB, config.Log, config.Validate, xenditPayoutRepository, config.XenditClient, walletRepository, userRepository)
	payoutUseCase := usecase.NewPayoutUseCase(config.DB, config.Log, config.Validate, payoutRepository, xenditPayoutUseCase, walletRepository, userRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	categoryController := http.NewCategoryController(categoryUseCase, config.Log)
	productController := http.NewProductController(productUseCase, config.Log)
	orderController := http.NewOrderController(orderUseCase, config.Log)
	discountCouponController := http.NewDiscountCouponController(discountCouponUseCase, config.Log)
	deliveryController := http.NewDeliveryController(deliveryUseCase, config.Log)
	productReviewController := http.NewProductReviewController(productReviewUseCase, config.Log)
	midtransSnapOrderController := http.NewMidtransSnapOrderController(midtransSnapOrderUseCase, config.Log)
	midtransCoreAPIOrderController := http.NewMidtransCoreAPIOrderController(midtransCoreApiOrderUseCase, config.Log)
	xenditQRCodeTransactionController := xenditController.NewXenditQRCodeTransctionController(xenditTransactionQRCodeUseCase, config.Log, config.DB)
	xenditCallbackController := xenditController.NewXenditCallbackController(xenditCallbackUseCase, config.Log)
	applicationController := http.NewApplicationController(applicationUseCase, config.Log)
	cartController := http.NewCartController(cartUseCase, config.Log)
	xenditPayoutController := xenditController.NewXenditPayoutController(xenditPayoutUseCase, config.Log, config.DB)
	payoutController := http.NewPayoutController(payoutUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	roleMiddleware := middleware.NewRole(userUseCase)
	authXenditMiddleware := middleware.NewAuthXenditCallback(*config.Config, config.Log)

	routeConfig := route.RouteConfig{
		App:                               config.App,
		UserController:                    userController,
		AddressController:                 addressController,
		CategoryController:                categoryController,
		ProductController:                 productController,
		OrderController:                   orderController,
		DiscountCouponController:          discountCouponController,
		DeliveryController:                deliveryController,
		ProductReviewController:           productReviewController,
		MidtransSnapOrderController:       midtransSnapOrderController,
		MidtransCoreAPIOrderController:    midtransCoreAPIOrderController,
		XenditQRCodeTransactionController: xenditQRCodeTransactionController,
		PayoutController:                  payoutController,
		XenditCallbackController:          xenditCallbackController,
		XenditPayoutController:            xenditPayoutController,
		ApplicationController:             applicationController,
		CartController:                    cartController,
		AuthMiddleware:                    authMiddleware,
		RoleMiddleware:                    roleMiddleware,
		AuthXenditMiddleware:              authXenditMiddleware,
	}
	routeConfig.Setup()
}
