// #region error-types
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

// #endregion

// #region map-status
// MapStatus は Go 1.26 の errors.AsType[T] で
// エラーを HTTP ステータスに対応付けます。
func MapStatus(err error) int {
	if _, ok := errors.AsType[*ValidationError](err); ok {
		return http.StatusBadRequest
	}
	if _, ok := errors.AsType[*NotFoundError](err); ok {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

// #endregion
