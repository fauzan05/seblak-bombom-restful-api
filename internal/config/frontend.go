package config

import (
	"seblak-bombom-restful-api/internal/model"

	"github.com/spf13/viper"
)

func NewFrontEndConfig(viper *viper.Viper) *model.FrontEndConfig {
	getBaseURL := viper.GetString("front_end.base_url")
	newFrontEnd := new(model.FrontEndConfig)
	newFrontEnd.BaseURL = getBaseURL
	return newFrontEnd
}
