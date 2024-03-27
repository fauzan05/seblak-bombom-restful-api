package config

import (
	"seblak-bombom-restful-api/internal/helper"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	helper.HandleErrorWithPanic(err)

	return config
}