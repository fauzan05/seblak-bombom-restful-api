package mailer

import (
	"fmt"
	"net/smtp"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"time"
)

type SMTPMailer struct {
	AuthEmail    string
	AuthPassword string
	SenderName   string
	Host         string
	Port         int
}

func (s *SMTPMailer) Send(mail model.Mail) error {
	body := mail.Template.String()
	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", s.SenderName, s.AuthEmail),
		"To":           strings.Join(mail.To, ","),
		"Subject":      mail.Subject,
		"Date":         time.Now().Format(time.RFC1123Z),
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	if len(mail.Cc) > 0 {
		headers["Cc"] = strings.Join(mail.Cc, ",")
	}

	var msg strings.Builder
	for k, v := range headers {
		if _, err := msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v)); err != nil {
			return fmt.Errorf("failed to write header %s: %w", k, err)
		}
	}
	msg.WriteString("\r\n" + body)

	auth := smtp.PlainAuth("", s.AuthEmail, s.AuthPassword, s.Host)
	recipients := append(mail.To, mail.Cc...)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.AuthEmail, recipients, []byte(msg.String()))
}