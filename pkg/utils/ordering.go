package utils

import (
	"net/http"
	"strings"

	"gopher-social-backend-server/pkg/constants"
)

func ParseOrderByQueryParam(r *http.Request, key, defaultValue string) string {
	paramStr := r.URL.Query().Get(key)
	if paramStr == "" {
		return defaultValue
	}

	parts := strings.Fields(paramStr)
	return parts[0]
}

func ParseDescQueryParam(r *http.Request, key, defaultValue string) string {
	paramStr := r.URL.Query().Get(key)
	if paramStr == "" {
		return defaultValue
	}

	return paramStr
}

func IsSQLInjection(queryPart string) bool {
	_, exists := constants.SQLKeywords[strings.ToUpper(queryPart)]
	return exists
}
