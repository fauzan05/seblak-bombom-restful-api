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
	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenRepository)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, userRepository, addressRepository, userUseCase)
	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)
	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App: config.App,
		UserController: userController,
		AddressController: addressController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}