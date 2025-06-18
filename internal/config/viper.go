package config

import (
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/helper"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()
	wd, err := os.Getwd()
    if err != nil {
        helper.HandleErrorWithPanic(err)
    }

    rootPath := filepath.Join(wd)

	config.SetConfigFile(".env")
    config.AddConfigPath(rootPath)
    config.AutomaticEnv()

    if err := config.ReadInConfig(); err != nil {
        helper.HandleErrorWithPanic(err)
    }

    return config

}
