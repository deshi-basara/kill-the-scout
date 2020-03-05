package notifier

import (
	"net/smtp"
	"net/textproto"
	"os"
	"strings"

	"github.com/deshi-basara/kill-the-scout/database"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// SendMail send a mail to a given address
func SendMail(to string, expose database.Expose) {
	var subject strings.Builder
	subject.WriteString("[" + expose.Rent + "] ")
	subject.WriteString(expose.Title + " (")
	subject.WriteString(expose.Size + "/")
	subject.WriteString(expose.Rooms + ")")

	var content strings.Builder
	content.WriteString(expose.Address + "\n")
	content.WriteString(expose.URL + "\n")
	content.WriteString(expose.CreatedAt.String())

	e := &email.Email{
		To:      []string{to},
		From:    "scout <no-reply@screensstudio.com>",
		Subject: subject.String(),
		Text:    []byte(content.String()),
		// HTML:    []byte(content.String()),
		Headers: textproto.MIMEHeader{},
	}

	// get needed smtp data
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	err := e.Send(smtpServer, smtp.PlainAuth(
		"", smtpUser, smtpPass, smtpHost,
	))
	if err != nil {
		log.Fatalf("SendMail.Fatal: %s", err)
	} else {
		log.Info("Mail sent")
	}
}
