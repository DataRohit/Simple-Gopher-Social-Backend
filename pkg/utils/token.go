package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ACTIVATION_MAIL_EXPIRATION = GetEnvAsDuration("ACTIVATION_MAIL_EXPIRATION", "30m")
var PASSWORD_RESET_EXPIRATION = GetEnvAsDuration("PASSWORD_RESET_EXPIRATION", "30m")
var JWT_SECRET = GetEnvAsByteArr("JWT_SECRET", "b82d4b46c665de2f8d506caf26f889c4d1b4d279a94fb99ef1f2d46992b034e5")
var JWT_EXPIRATION = GetEnvAsDuration("JWT_EXPIRATION", "6h")

func GenerateActivationToken(userID string) string {
	expirationTime := time.Now().Add(ACTIVATION_MAIL_EXPIRATION)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWT_SECRET)
	return tokenString
}

func VerifyActivationToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Subject, nil
}

func GenerateAccessToken(userID string) (string, time.Time) {
	expirationTime := time.Now().Add(JWT_EXPIRATION)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWT_SECRET)
	return tokenString, expirationTime
}

func ValidateAccessToken(tokenStr string) bool {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	return claims.ExpiresAt.Time.After(time.Now())
}

func GeneratePasswordResetToken(userID string) string {
	expirationTime := time.Now().Add(PASSWORD_RESET_EXPIRATION)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWT_SECRET)
	return tokenString
}

func VerifyPasswordResetToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Subject, nil
}

func ParseAccessToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
