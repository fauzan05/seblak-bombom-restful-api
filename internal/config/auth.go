package config

import (
	"seblak-bombom-restful-api/internal/model"

	"github.com/spf13/viper"
)

func NewAuthConfig(viper *viper.Viper) *model.AuthConfig {
	adminCreationKey := viper.GetString("ADMIN_CREATION_KEY")
	newAuthConfig := new(model.AuthConfig)
	newAuthConfig.AdminCreationKey = adminCreationKey
	return newAuthConfig
}
