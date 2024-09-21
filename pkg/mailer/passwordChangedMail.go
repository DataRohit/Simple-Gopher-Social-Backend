package mailer

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

var password_changed_email_template = template.Must(template.ParseFiles("templates/authentication/password_changed_email.gtpl"))

type PasswordChangedEmailData struct {
	Email string
}

func SendPasswordChangedEmail(email string) {
	data := PasswordChangedEmailData{
		Email: email,
	}

	var body bytes.Buffer
	if err := password_changed_email_template.Execute(&body, data); err != nil {
		log.Printf("error rendering password changed email template: %v", err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Password Changed")
	mail.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("mailpit", 1025, "", "")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("error sending email: %v", err)
		return
	}

	log.Printf("password changed email sent to %s", email)
}
