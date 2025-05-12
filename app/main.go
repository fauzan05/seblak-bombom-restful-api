package main

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	// db := config.NewDatabaseProd(viperConfig, log) //prod
	// db := config.NewDatabaseTest(viperConfig, log) // test
	db := config.NewDatabaseDev(viperConfig, log) // dev
	// snapClient := config.NewMidtransSanboxSnapClient(viperConfig, log)
	// coreAPIClient := config.NewMidtransSanboxCoreAPIClient(viperConfig, log)
	xenditClient := config.NewXenditTestTransactions(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	email := config.NewSMTPMailerTest(viperConfig)
	pdf := config.NewPDFGenerator(log)
	app := config.NewFiber(viperConfig)
	
	// cors setting
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://seblak-bombom-api-consumer-app, http://localhost:8000, http://127.0.0.1:5500",  // Frontend are allowed (port 8000), if you use docker so you have to list the container name of api consumer (seblak-bombom-api-consumer)
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE",                                     // HTTP method are allowed
		AllowHeaders:     "Origin, Content-Type, X-Requested-With, Accept, Authorization", // Header are allowed
		AllowCredentials: true,
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		// SnapClient:    snapClient,
		// CoreAPIClient: coreAPIClient,
		XenditClient: xenditClient,
		Email:        email,
		PDF:          pdf,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
