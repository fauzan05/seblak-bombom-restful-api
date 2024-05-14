package main

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabaseProd(viperConfig, log) //prod
	// db := config.NewDatabaseTest(viperConfig, log) // test
	// db := config.NewDatabaseDev(viperConfig, log) // dev
	snapClient := config.NewMidtransSanboxSnapClient(viperConfig, log)
	coreAPIClient := config.NewMidtransSanboxCoreAPIClient(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

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
