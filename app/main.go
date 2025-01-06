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
	snapClient := config.NewMidtransSanboxSnapClient(viperConfig, log)
	coreAPIClient := config.NewMidtransSanboxCoreAPIClient(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	
	// cors setting
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "http://localhost:9000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:            db,
		App:           app,
		Log:           log,
		Validate:      validate,
		Config:        viperConfig,
		SnapClient:    snapClient,
		CoreAPIClient: coreAPIClient,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
