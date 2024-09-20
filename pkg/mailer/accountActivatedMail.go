package mailer

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

var accountActivatedEmailTemplate = template.Must(template.ParseFiles("templates/authentication/account_activated_email.tmpl"))

type AccountActivatedEmailData struct {
	Email string
}

func SendAccountActivatedEmail(email string) {
	data := AccountActivatedEmailData{
		Email: email,
	}

	var body bytes.Buffer
	if err := accountActivatedEmailTemplate.Execute(&body, data); err != nil {
		log.Printf("error rendering account activated email template: %v", err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Account Activated")
	mail.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("mailpit", 1025, "", "")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("error sending email: %v", err)
		return
	}

	log.Printf("account activated confirmation email sent to %s", email)
}
