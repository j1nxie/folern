package models

import "fmt"

type FolernResponse[T any] struct {
	Success bool   `json:"success"`
	Body    T      `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
	Cat     string `json:"cat"`
}

type FolernError struct {
	Message string `json:"message"`
}

func (e FolernError) Error() string {
	return e.Message
}

func SuccessResponse[T any](statusCode int, body T) FolernResponse[T] {
	return FolernResponse[T]{
		Success: true,
		Body:    body,
		Cat:     fmt.Sprintf("https://http.cat/%d", statusCode),
	}
}

func ErrorResponse[T any](statusCode int, err error) FolernResponse[T] {
	return FolernResponse[T]{
		Success: false,
		Error:   err,
		Cat:     fmt.Sprintf("https://http.cat/%d", statusCode),
	}
}
