package utils

import (
	"fmt"
	"time"
)

var ACTIVATION_MAIL_EXPIRATION = GetEnvAsDuration("ACTIVATION_MAIL_EXPIRATION", "30m")

func SendActivationEmail(email, token string) {
	expiration := time.Now().Add(ACTIVATION_MAIL_EXPIRATION).Format(time.RFC1123)
	message := fmt.Sprintf("Please activate your account by visiting the following link: http://localhost:8080/activate/%s. The link expires at %s", token, expiration)

	fmt.Printf("Sending email to %s: %s", email, message)
}
