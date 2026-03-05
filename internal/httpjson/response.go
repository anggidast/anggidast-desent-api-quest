package httpjson

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"desent-api-quest/internal/domain"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func WriteError(w http.ResponseWriter, err error) {
	status, body := MapError(err)
	WriteJSON(w, status, body)
}

func MapError(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()}
	case errors.Is(err, domain.ErrBadRequest):
		return http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: err.Error()}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: err.Error()}
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()}
	default:
		message := "internal server error"
		if trimmed := strings.TrimSpace(err.Error()); trimmed != "" && !errors.Is(err, domain.ErrInternal) {
			message = trimmed
		}
		return http.StatusInternalServerError, ErrorResponse{Error: "internal_error", Message: message}
	}
}
