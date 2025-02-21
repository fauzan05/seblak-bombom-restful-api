package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"seblak-bombom-restful-api/internal/config"
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
	app := config.NewFiber(viperConfig)

	// cors setting
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8000",                                         // Frontend yang diizinkan (port 8000)
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE",                                     // Metode HTTP yang diizinkan
		AllowHeaders:     "Origin, Content-Type, X-Requested-With, Accept, Authorization", // Header yang diizinkan
		AllowCredentials: true,                                                            // Mengizinkan pengiriman cookie
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
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
