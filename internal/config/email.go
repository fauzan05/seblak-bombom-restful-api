package config

import (
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"

	"github.com/spf13/viper"
)

func NewEmailWorker(viper *viper.Viper) *mailer.EmailWorker {
	smtpMailer := &mailer.SMTPMailer{
		AuthEmail:    viper.GetString("EMAIL_TEST_EMAIL"),
		AuthPassword: viper.GetString("EMAIL_TEST_PASSWORD"),
		SenderName:   viper.GetString("EMAIL_TEST_SENDER_NAME"),
		Host:         viper.GetString("EMAIL_TEST_HOST"),
		Port:         viper.GetInt("EMAIL_TEST_PORT"),
	}

	worker := &mailer.EmailWorker{
		Mailer:     smtpMailer,
		MailQueue:  make(chan model.Mail, 1000), // buffer besar
		WorkerSize: 10,                          // atau configurable
	}

	worker.Start()
	return worker
}
