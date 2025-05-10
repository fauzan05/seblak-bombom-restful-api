package tests

import (
	"seblak-bombom-restful-api/internal/config"
	"seblak-bombom-restful-api/internal/helper/mailer"

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

var email *mailer.SMTPMailer

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabaseTest(viperConfig, log)
	email = config.NewSMTPMailerTest(viperConfig)
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Email:    email,
	})
}
