package domain

import "errors"

var (
	ErrBadRequest   = errors.New("bad request")
	ErrValidation   = errors.New("validation error")
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
)

type FieldError struct {
	Code    error
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}

func (e *FieldError) Unwrap() error {
	return e.Code
}
