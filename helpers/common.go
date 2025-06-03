package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/mrgThang/flashcard-be/dto"
)

func WriteJSONError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := dto.ApiResponse[any]{
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
	response := dto.ApiResponse[T]{
		Code:    code,
		Message: "Success",
		Data:    data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
