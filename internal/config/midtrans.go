package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


func NewMidtransSanbox(viper *viper.Viper, log *logrus.Logger) *snap.Client {
	var client snap.Client
	serverKey := viper.GetString("midtrans.sandbox.server_key")
	client.New(serverKey, midtrans.Sandbox)

	return &client
}