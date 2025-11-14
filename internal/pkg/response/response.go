package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, success bool, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Success: success, Message: message, Data: data})
}

func Success(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusOK, true, message, data)
}

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, false, message, nil)
}

func InternalError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, false, message, nil)
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, false, message, nil)
}
