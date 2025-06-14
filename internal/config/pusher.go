package config

import (
	"github.com/pusher/pusher-http-go/v5"
	"github.com/spf13/viper"
)

func NewPusherClient(viper *viper.Viper) pusher.Client {
	appId := viper.GetString("pusher.app_id")
	key := viper.GetString("pusher.key")
	secret := viper.GetString("pusher.secret")
	cluster := viper.GetString("pusher.cluster")
	secure := viper.GetBool("pusher.secure")
	pusherClient := pusher.Client{
		AppID:   appId,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  secure,
	}

	return pusherClient
}
