package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/httpjson"
)

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrBadRequest, Message: "invalid JSON body"})
		return
	}

	var payload any
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&payload); err != nil {
		httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrBadRequest, Message: "invalid JSON body"})
		return
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrBadRequest, Message: "request body must contain a single JSON value"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
