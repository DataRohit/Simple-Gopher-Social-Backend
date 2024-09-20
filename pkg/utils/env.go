package utils

import (
	"os"
	"strconv"
	"time"
)

func GetEnvAsString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetEnvAsInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

func GetEnvAsBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return boolValue
}

func GetEnvAsDuration(key string, fallback string) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		duration, _ := time.ParseDuration(fallback)
		return duration
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		duration, _ = time.ParseDuration(fallback)
	}

	return duration
}

func GetEnvAsByteArr(key string, fallback string) []byte {
	value, exists := os.LookupEnv(key)
	if !exists {
		return []byte(fallback)
	}
	return []byte(value)
}
