package httpjson

import (
	"encoding/json"
	"errors"
	"io"
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

func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return &domain.FieldError{Code: domain.ErrBadRequest, Message: "request body is required"}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return &domain.FieldError{Code: domain.ErrBadRequest, Message: "request body is required"}
		}
		return &domain.FieldError{Code: domain.ErrBadRequest, Message: "invalid JSON body"}
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return &domain.FieldError{Code: domain.ErrBadRequest, Message: "request body must contain a single JSON value"}
	}

	return nil
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
