package config

import (
	"seblak-bombom-restful-api/internal/helper/mailer"
	"seblak-bombom-restful-api/internal/model"

	"github.com/spf13/viper"
)

func NewEmailWorker(viper *viper.Viper) *mailer.EmailWorker {
	smtpMailer := &mailer.SMTPMailer{
		AuthEmail:    viper.GetString("email.test.email"),
		AuthPassword: viper.GetString("email.test.password"),
		SenderName:   viper.GetString("email.test.sender_name"),
		Host:         viper.GetString("email.test.host"),
		Port:         viper.GetInt("email.test.port"),
	}

	worker := &mailer.EmailWorker{
		Mailer:     smtpMailer,
		MailQueue:  make(chan model.Mail, 1000), // buffer besar
		WorkerSize: 10,                          // atau configurable
	}

	worker.Start()
	return worker
}
