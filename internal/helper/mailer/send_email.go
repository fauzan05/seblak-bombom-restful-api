package mailer

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"mime/quotedprintable"
	"net/smtp"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type EmailWorker struct {
	Mailer     *SMTPMailer
	MailQueue  chan model.Mail
	WorkerSize int
	wg         sync.WaitGroup
}

type SMTPMailer struct {
	AuthEmail    string
	AuthPassword string
	SenderName   string
	Host         string
	Port         int
}

func (s *SMTPMailer) Send(mail model.Mail) error {
	boundary := helper_others.GenerateBoundary() // batas antara bagian HTML dan attachment

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", s.SenderName, s.AuthEmail),
		"To":           strings.Join(mail.To, ","),
		"Subject":      mail.Subject,
		"Date":         time.Now().Format(time.RFC1123Z),
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary),
	}
	if len(mail.Cc) > 0 {
		headers["Cc"] = strings.Join(mail.Cc, ",")
	}

	var msg strings.Builder

	// tulis headers
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	// bagian HTML
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qp := quotedprintable.NewWriter(&msg)
	qp.Write([]byte(mail.Template.String()))
	qp.Close()

	// bagian attachment
	for _, att := range mail.Attachments {
		msg.WriteString("\r\n--" + boundary + "\r\n")
		msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", att.MimeType, att.Filename))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.Filename))

		b64 := make([]byte, base64.StdEncoding.EncodedLen(len(att.Content)))
		base64.StdEncoding.Encode(b64, att.Content)

		// tulis base64 dengan newline setiap 76 karakter (standar MIME)
		for i := 0; i < len(b64); i += 76 {
			end := i + 76
			if end > len(b64) {
				end = len(b64)
			}
			msg.Write(b64[i:end])
			msg.WriteString("\r\n")
		}
	}

	msg.WriteString("\r\n--" + boundary + "--")

	auth := smtp.PlainAuth("", s.AuthEmail, s.AuthPassword, s.Host)
	recipients := append(mail.To, mail.Cc...)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.AuthEmail, recipients, []byte(msg.String()))
}

func (w *EmailWorker) Start() {
	for i := 0; i < w.WorkerSize; i++ {
		w.wg.Add(1)
		go w.runWorker(i)
	}
}

func (w *EmailWorker) runWorker(id int) {
	defer w.wg.Done()
	for mail := range w.MailQueue {
		err := w.Mailer.Send(mail)
		if err != nil {
			log.Printf("[Worker %d] Failed send email to %v: %v", id, mail.To, err)
		} else {
			log.Printf("[Worker %d] Email successfully sent to %v", id, mail.To)
		}
	}
	log.Printf("[Worker %d] Finished", id)
}

func (w *EmailWorker) Stop() {
	close(w.MailQueue) // Ini penting untuk keluarin semua goroutine dari loop
	w.wg.Wait()        // Tunggu semua selesai
}

func (c *EmailWorker)SendEmail(log *logrus.Logger, to []string, cc []string, subject string, baseTemplatePath string, childTemplatePath string, data map[string]any) error {
	tmpl, err := template.ParseFiles(baseTemplatePath, childTemplatePath)
	if err != nil {
		log.Warnf("failed to parse template file html : %+v", err)
		return err
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.ExecuteTemplate(bodyBuilder, "base", data)
	if err != nil {
		log.Warnf("failed to execute template file html : %+v", err)
		return err
	}

	newMail := new(model.Mail)
	newMail.To = to
	if len(cc) > 0 {
		newMail.Cc = cc
	}
	newMail.Subject = subject
	newMail.Template = *bodyBuilder
	c.Mailer.SenderName = fmt.Sprintf("System %s", data["CompanyTitle"])

	select {
	case c.MailQueue <- *newMail:
		return nil
	default:
		return fmt.Errorf("email queue full, failed to send to %s", strings.Join(to, ", "))
	}
}
