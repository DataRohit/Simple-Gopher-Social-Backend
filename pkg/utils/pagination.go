package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func ParseLimitOffsetQueryParam(r *http.Request, key string, defaultValue, min, max int) (int, error) {
	paramStr := r.URL.Query().Get(key)
	if paramStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil || value < min || (max >= 0 && value > max) {
		return 0, fmt.Errorf("invalid %s: must be between %d and %d", key, min, max)
	}

	return value, nil
}
