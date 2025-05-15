package tests

import (
	"fmt"
	"html/template"
	"seblak-bombom-restful-api/internal/interfaces"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	newEmail := new(model.Mail)
	newEmail.To = []string{"F3196813@gmail.com"}
	newEmail.Cc = []string{}
	newEmail.Subject = "Testing Email Ya"
	templatePath := "../tests/templates/email/testing_email.html"
	tmpl, err := template.ParseFiles(templatePath)
	assert.Nil(t, err)

	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]string{
		"Message": "Ini adalah testing email ya",
		"Year":    time.Now().Format("2006"),
		"Sender":  "Fauzan Nur Hidayat",
	})
	assert.Nil(t, err)

	newEmail.Template = *bodyBuilder

	select {
	case email.MailQueue <- *newEmail:
		fmt.Println("Email masuk ke queue")
	case <-time.After(2 * time.Second):
		t.Fatal("Gagal enqueue email: timeout")
	}

	time.Sleep(3 * time.Second)
	assert.Nil(t, err)
}

type JustPrintSubject struct {
	Name string
}

func (s *JustPrintSubject) Send(mail model.Mail) error {
	text := fmt.Sprintf("NAMA SUBJECT : %s", mail.Subject)
	fmt.Println(text)
	return nil
}

type MergeAllMailer struct {
	MailerTask interfaces.Mailer
}

func (n *MergeAllMailer) Notify(mail model.Mail) {
	_ = n.MailerTask.Send(mail)
}

// func TestSendEmailAndLearnInterface(t *testing.T) {
// 	newEmail := new(model.Mail)
// 	newEmail.To = []string{"F3196813@gmail.com"}
// 	newEmail.Cc = []string{}
// 	newEmail.Subject = "Testing Email Ya"
// 	templatePath := "../internal/helper/templates/testing_email.html"
// 	tmpl, err := template.ParseFiles(templatePath)
// 	assert.Nil(t, err)

// 	bodyBuilder := new(strings.Builder)
// 	err = tmpl.Execute(bodyBuilder, map[string]string{
// 		"Message": "Ini adalah testing email ya",
// 		"Year":    time.Now().Format("2006"),
// 		"Sender":  "Fauzan Nur Hidayat",
// 	})
// 	assert.Nil(t, err)

// 	// ini akan ngirim email
// 	newEmail.Template = *bodyBuilder
// 	merge := MergeAllMailer{
// 		MailerTask: email,
// 	}

// 	err = merge.MailerTask.Send(*newEmail)
// 	assert.Nil(t, err)

// 	// ini akan ngeprint aja
// 	var newSubject JustPrintSubject
// 	merge = MergeAllMailer{
// 		MailerTask: &newSubject,
// 	}
// 	err = merge.MailerTask.Send(*newEmail)
// 	assert.Nil(t, err)
// }
