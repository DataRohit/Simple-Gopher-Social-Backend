package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to write JSON response: %v", err)
	}
	return nil
}

func WriteError(w http.ResponseWriter, status int, err string) {
	errorResponse := map[string]string{"error": err}

	if writeErr := WriteJSON(w, status, errorResponse); writeErr != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
	}
}
