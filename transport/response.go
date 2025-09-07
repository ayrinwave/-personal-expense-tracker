package transport

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — это стандартная структура для JSON-ответов с ошибкой.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError отправляет JSON-ответ с ошибкой и заданным HTTP-статусом.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

// RespondWithJSON отправляет JSON-ответ с любыми данными и заданным HTTP-статусом.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
