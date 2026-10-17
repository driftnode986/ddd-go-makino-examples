// #region error-mapper
package presentation

import (
	"errors"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

func mapErrorToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidStatus):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// #endregion
