package errors

import (
	stderrors "errors"
	"net/http"
)

type Kind string

const (
	KindValidation Kind = "validation"
	KindNotFound   Kind = "not_found"
	KindConflict   Kind = "conflict"
	KindInternal   Kind = "internal"
)

type AppError struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "internal server error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(kind Kind, message string, err error) error {
	return &AppError{
		Kind:    kind,
		Message: message,
		Err:     err,
	}
}

func KindOf(err error) Kind {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Kind
	}
	return KindInternal
}

func MessageOf(err error) string {
	if err == nil {
		return ""
	}
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Error()
	}
	return "internal server error"
}

func HTTPStatus(err error) int {
	switch KindOf(err) {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
