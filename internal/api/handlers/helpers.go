package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse отправляет стандартный ответ об ошибке.
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "error",
		"code":   status,
		"errors": message,
	})
}
