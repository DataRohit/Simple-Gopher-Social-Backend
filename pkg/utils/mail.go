package utils

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

var ACTIVATION_MAIL_EXPIRATION = GetEnvAsDuration("ACTIVATION_MAIL_EXPIRATION", "30m")

func SendActivationEmail(email, token string) {
	expiration := time.Now().Add(1 * time.Hour).Format(time.RFC1123)
	subject := "Activate Your Account"
	body := fmt.Sprintf("Please activate your account by visiting the following link: http://localhost:8080/activate/%s. The link expires at %s", token, expiration)

	mail := gomail.NewMessage()
	mail.SetHeader("From", "no-reply@gopher.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	dialer := gomail.NewDialer("mailpit", 1025, "", "")

	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}

	log.Printf("Activation email sent to %s", email)
}
