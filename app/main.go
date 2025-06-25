package main

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	env := viperConfig.GetString("ENV")
	if env == "" {
		log.Fatal("ENV is not set")
	}

	if env != "prod" && env != "test" && env != "dev" && env != "docker" {
		log.Fatalf("Invalid ENV value: %s. Must be one of: prod, test, dev, docker", env)
	}

	var db *gorm.DB // Use an empty interface to hold different database types
	if env == "dev" {
		fmt.Println("Running in development mode")
		db = config.NewDatabaseDev(viperConfig, log) // dev
	}

	if env == "prod" {
		fmt.Println("Running in production mode")
		db = config.NewDatabaseProd(viperConfig, log) //prod
	}

	if env == "test" {
		fmt.Println("Running in test mode")
		db = config.NewDatabaseTest(viperConfig, log) // test
	}

	if env == "docker" {
		fmt.Println("Running in docker mode")
		db = config.NewDatabaseDocker(viperConfig, log) // docker
	}

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
				// "https://seblak-bombom-api-consumer-production.up.railway.app": true,
				"http://localhost:3000": true,
				"https://seblak.fznh-dev.my.id": true,
			}
			return allowed[origin]
		},
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		// ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Authorization,Set-Cookie",
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

	webPort := viperConfig.GetInt("WEB_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}

	defer email.Stop()
}
