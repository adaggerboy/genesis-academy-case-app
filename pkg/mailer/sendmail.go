package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"github.com/adaggerboy/genesis-academy-case-app/config"
	confModel "github.com/adaggerboy/genesis-academy-case-app/models/config"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/3rd/openexchangeapi"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/database"
)

// go:embed template.html
var templateData string

type EmailData struct {
	Rate float32
}

func SendMail(to string, subject string, data EmailData, conf confModel.SMTPConfig) error {

	tmpl, err := template.New("emailTemplate").Parse(templateData)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", conf.User, conf.Password, conf.Password)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", conf.User, to, subject)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", conf.Host, conf.Port), auth, conf.User, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func GoThroughSubscriptions() error {
	subs, err := database.GetDatabase().GetSubscriptions()
	if err != nil {
		return err
	}
	for _, e := range subs {
		go func() {
			rate, err := openexchangeapi.RequestUSDPairCached("UAH")
			if err != nil {
				log.Printf("Error: getting actual rates: %s", err)
			}
			SendMail(e, "Current USD/UAH rate", EmailData{
				Rate: rate,
			}, config.GlobalConfig.SMTPConfig)
			if err != nil {
				log.Printf("Error: sending email: %s", err)
			}
		}()
	}
	return nil
}
