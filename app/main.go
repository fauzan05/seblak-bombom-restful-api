package main

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabaseProd(viperConfig, log) //prod
	// db := config.NewDatabaseTest(viperConfig, log) // test
	// db := config.NewDatabaseDev(viperConfig, log) // dev
	// db := config.NewDatabaseDocker(viperConfig, log)
	xenditClient := config.NewXenditTestTransactions(viperConfig, log)
	validate := config.NewValidator()
	email := config.NewEmailWorker(viperConfig)
	authConfig := config.NewAuthConfig(viperConfig)
	frontEndConfig := config.NewFrontEndConfig(viperConfig)
	pdf := config.NewPDFGenerator(log)
	pusherClient := config.NewPusherClient(viperConfig)

	app := config.NewFiber(viperConfig)

	// cors setting
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			allowed := map[string]bool{
				"http://localhost:3000":                                         true,
				"http://seblak-bombom-api-consumer-app":                         true,
				"https://seblak-bombom-api-consumer-production.up.railway.app": true,
			}
			return allowed[origin]
		},
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:             db,
		App:            app,
		Log:            log,
		Validate:       validate,
		Config:         viperConfig,
		XenditClient:   xenditClient,
		Email:          email,
		PDF:            pdf,
		AuthConfig:     authConfig,
		FrontEndConfig: frontEndConfig,
		PusherClient:   pusherClient,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}

	defer email.Stop()
}
