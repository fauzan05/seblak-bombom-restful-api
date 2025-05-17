package config

import (
	"seblak-bombom-restful-api/internal/model"

	"github.com/spf13/viper"
)

func NewAuthConfig(viper *viper.Viper) *model.AuthConfig {
	adminCreationKey := viper.GetString("auth_admin.admin_creation_key")
	newAuthConfig := new(model.AuthConfig)
	newAuthConfig.AdminCreationKey = adminCreationKey
	return newAuthConfig
}
