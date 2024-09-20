package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func ReadJSON(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	const maxPayloadSize = 1 * 1024 * 1024
	buffer := bytes.NewBuffer(make([]byte, 0, maxPayloadSize))
	lr := io.LimitedReader{R: r.Body, N: maxPayloadSize}

	_, err := io.Copy(buffer, &lr)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("payload size exceeds the maximum allowed size of 1 MB")
		}
		return nil, fmt.Errorf("unable to read request body: %v", err)
	}

	decoder := json.NewDecoder(buffer)
	decoder.DisallowUnknownFields()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("invalid JSON format or unknown fields: %v", err)
	}

	return buffer.Bytes(), nil
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
