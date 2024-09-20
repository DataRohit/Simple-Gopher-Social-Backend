package mailer

import (
	"bytes"
	"gopher-social-backend-server/pkg/constants"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

var oauth_welcome_email_template = template.Must(template.ParseFiles("templates/authentication/oauth_welcome_email.tmpl"))

type OAuthWelcomeEmailData struct {
	Email    string
	Provider constants.OAuthProvider
}

func SendOAuthWelcomeEmail(email string, provider constants.OAuthProvider) {
	data := OAuthWelcomeEmailData{
		Email:    email,
		Provider: provider,
	}

	var body bytes.Buffer
	if err := oauth_welcome_email_template.Execute(&body, data); err != nil {
		log.Printf("error rendering oAuth welcome email template: %v", err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "User Registered")
	mail.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("mailpit", 1025, "", "")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("error sending email: %v", err)
		return
	}

	log.Printf("welcome email sent to %s", email)
}
