package httpjson

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"desent-api-quest/internal/domain"
)

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