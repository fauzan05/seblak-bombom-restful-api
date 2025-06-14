package tests

import (
	"os"
	"seblak-bombom-restful-api/internal/config"
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

var email *mailer.EmailWorker

var authConfig *model.AuthConfig

var frontEndConfig *model.FrontEndConfig

func init() {
	os.Setenv("TZ", "UTC")
	time.Local = time.UTC // ini yang benar-benar bikin time.Now() jadi UTC
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator()
	app = config.NewFiber(viperConfig)
	db = config.NewDatabaseDockerTest(viperConfig, log)
	email = config.NewEmailWorker(viperConfig)
	authConfig = config.NewAuthConfig(viperConfig)
	frontEndConfig = config.NewFrontEndConfig(viperConfig)
	pusherClient := config.NewPusherClient(viperConfig)
	config.Bootstrap(&config.BootstrapConfig{
		DB:             db,
		App:            app,
		Log:            log,
		Validate:       validate,
		Config:         viperConfig,
		Email:          email,
		AuthConfig:     authConfig,
		FrontEndConfig: frontEndConfig,
		PusherClient:   pusherClient,
	})
}
