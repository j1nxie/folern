package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FolernResponse[T any] struct {
	Success     bool   `json:"success"`
	Description string `json:"description,omitempty"`
	Body        *T     `json:"body,omitempty"`
	Cat         string `json:"cat"`
}

func SuccessResponse[T any](w http.ResponseWriter, statusCode int, description string, body T) {
	resp := FolernResponse[T]{
		Success:     true,
		Description: description,
		Body:        &body,
		Cat:         fmt.Sprintf("https://http.cat/%d", statusCode),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(resp)
}

func ErrorResponse[T any](w http.ResponseWriter, statusCode int, description string) {
	resp := FolernResponse[T]{
		Success:     false,
		Description: description,
		Cat:         fmt.Sprintf("https://http.cat/%d", statusCode),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(resp)
}
