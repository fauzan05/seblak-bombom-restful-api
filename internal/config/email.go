package config

import (
	"seblak-bombom-restful-api/internal/helper/mailer"

	"github.com/spf13/viper"
)

func NewSMTPMailerTest(viper *viper.Viper) *mailer.SMTPMailer {
	return &mailer.SMTPMailer{
		AuthEmail:    viper.GetString("email.test.email"),
		AuthPassword: viper.GetString("email.test.password"),
		SenderName:   viper.GetString("email.test.sender_name"),
		Host:         viper.GetString("email.test.host"),
		Port:         viper.GetInt("email.test.port"),
	}
}
