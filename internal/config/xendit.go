package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xendit/xendit-go/v6"
)

func NewXenditTestTransactions(viper *viper.Viper, log *logrus.Logger) *xendit.APIClient {
	apiKey := viper.GetString("xendit.test.api_key")
	xenditClient := xendit.NewClient(apiKey)
	return xenditClient
}
