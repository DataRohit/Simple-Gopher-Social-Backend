package mailer

import (
	"bytes"
	"gopher-social-backend-server/pkg/utils"
	"html/template"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

var PASSWORD_RESET_EXPIRATION = utils.GetEnvAsDuration("PASSWORD_RESET_EXPIRATION", "30m")

var password_reset_email_template = template.Must(template.ParseFiles("templates/authentication/password_reset_email.gtpl"))

type PasswordResetEmailData struct {
	Email      string
	Token      string
	Expiration string
}

func SendPasswordResetEmail(email, token string) {
	expiration := time.Now().Add(PASSWORD_RESET_EXPIRATION).Format(time.RFC1123)

	data := PasswordResetEmailData{
		Email:      email,
		Token:      token,
		Expiration: expiration,
	}

	var body bytes.Buffer
	if err := password_reset_email_template.Execute(&body, data); err != nil {
		log.Printf("error rendering password reset email template: %v", err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Reset Your Password")
	mail.SetBody("text/html", body.String())

	dialer := gomail.NewDialer("mailpit", 1025, "", "")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("error sending email: %v", err)
		return
	}

	log.Printf("password reset email sent to %s", email)
}
