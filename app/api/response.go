package api

import (
	"encoding/json"
	"net/http"
)

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`

	//Optional fields
	Page    int  `json:"page,omitempty"`
	Limit   int  `json:"limit,omitempty"`
	Count   int  `json:"count,omitempty"`
	HasNext bool `json:"hasNext,omitempty"`
	Total   int  `json:"total,omitempty"`
}

func OKResponse[T any](w http.ResponseWriter, statusCode int, response ApiResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ApiResponse[any]{
		Success: false,
		Message: message,
		Data:    nil,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
