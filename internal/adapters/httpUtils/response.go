package httputil

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amangirdhar210/meeting-room/internal/core/domain"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func HandleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		RespondWithError(w, http.StatusNotFound, "resource not found")
	case domain.ErrInvalidInput:
		RespondWithError(w, http.StatusBadRequest, "invalid input data")
	case domain.ErrUnauthorized:
		RespondWithError(w, http.StatusUnauthorized, "unauthorized access")
	case domain.ErrConflict:
		RespondWithError(w, http.StatusConflict, "resource conflict")
	case domain.ErrRoomUnavailable:
		RespondWithError(w, http.StatusConflict, "room not available for the selected time slot")
	case domain.ErrTimeRangeInvalid:
		RespondWithError(w, http.StatusBadRequest, "invalid start or end time for booking")
	default:
		log.Printf("Unhandled error: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "internal server error")
	}
}
