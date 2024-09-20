package mailer

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

var github_welcome_email_template = template.Must(template.ParseFiles("templates/authentication/github_welcome_email.tmpl"))

type GithubWelcomeEmailData struct {
	Email string
}

func SendGitHubWelcomeEmail(email string) {
	data := GithubWelcomeEmailData{
		Email: email,
	}

	var body bytes.Buffer
	if err := github_welcome_email_template.Execute(&body, data); err != nil {
		log.Printf("error rendering github welcome email template: %v", err)
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
