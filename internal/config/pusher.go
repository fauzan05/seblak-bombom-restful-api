package config

import (
	"github.com/pusher/pusher-http-go/v5"
	"github.com/spf13/viper"
)

func NewPusherClient(viper *viper.Viper) pusher.Client {
	appId := viper.GetString("PUSHER_APP_ID")
	key := viper.GetString("PUSHER_KEY")
	secret := viper.GetString("PUSHER_SECRET")
	cluster := viper.GetString("PUSHER_CLUSTER")
	secure := viper.GetBool("PUSHER_SECURE")
	pusherClient := pusher.Client{
		AppID:   appId,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  secure,
	}

	return pusherClient
}
