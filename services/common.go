package services

import (
	"encoding/json"
	"net/http"
)

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func WriteJSONError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := ApiResponse[any]{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func WriteJSONResponse[T any](w http.ResponseWriter, code int, data T) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := ApiResponse[T]{
		Code:    code,
		Message: "Success",
		Data:    data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
