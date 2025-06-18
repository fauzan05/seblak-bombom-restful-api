package config

import (
	"os"
	// "path/filepath"
	// "seblak-bombom-restful-api/internal/helper"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
    config := viper.New()

    config.AutomaticEnv()

    if _, err := os.Stat(".env"); err == nil {
        config.SetConfigFile(".env")
        config.SetConfigType("env")
        _ = config.ReadInConfig()
    }

    return config
}

