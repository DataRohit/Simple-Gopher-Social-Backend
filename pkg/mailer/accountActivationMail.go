package mailer

import (
	"bytes"
	"gopher-social-backend-server/pkg/utils"
	"html/template"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

var ACTIVATION_MAIL_EXPIRATION = utils.GetEnvAsDuration("ACTIVATION_MAIL_EXPIRATION", "30m")

var activationEmailTemplate = template.Must(template.ParseFiles("templates/authentication/account_activation_email.tmpl"))

type ActivationEmailData struct {
	Email      string
	Token      string
	Expiration string
}

func SendActivationEmail(email, token string) {
	expiration := time.Now().Add(ACTIVATION_MAIL_EXPIRATION).Format(time.RFC1123)

	data := ActivationEmailData{
		Email:      email,
		Token:      token,
		Expiration: expiration,
	}

	var body bytes.Buffer
	if err := activationEmailTemplate.Execute(&body, data); err != nil {
		log.Printf("error rendering activation email template: %v", err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Activate Your Account")
	mail.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("mailpit", 1025, "", "")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("error sending email: %v", err)
		return
	}

	log.Printf("activation email sent to %s", email)
}
