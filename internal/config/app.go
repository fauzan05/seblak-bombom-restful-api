package config

import (
	"seblak-bombom-restful-api/internal/delivery/http"
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

	// setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup middleware

	routeConfig := route.RouteConfig{
		App: config.App,
		UserController: userController,
	}
	routeConfig.Setup()
	
}