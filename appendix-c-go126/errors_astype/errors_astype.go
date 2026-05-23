// Package errors_astype demonstrates Go 1.26's errors.AsType[T].
package errors_astype

import (
	"errors"
	"fmt"
	"net/http"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s id=%s", e.Resource, e.ID)
}

// MapStatus uses Go 1.26 errors.AsType[T] to map domain errors to HTTP status.
func MapStatus(err error) int {
	if _, ok := errors.AsType[*ValidationError](err); ok {
		return http.StatusBadRequest
	}
	if _, ok := errors.AsType[*NotFoundError](err); ok {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
