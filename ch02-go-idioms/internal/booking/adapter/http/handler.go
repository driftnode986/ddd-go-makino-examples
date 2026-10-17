// #region handler-imports
package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// #endregion

// #region handle-create-booking-a
// HTTP のアダプタでドメインエラーを HTTP ステータスに対応させる
func (h *Handler) handleCreateBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := h.useCase.Execute(r.Context(), buildInput(r))
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}

	// #endregion
	// #region handle-create-booking-b
	// このユースケースが返すエラーだけを並べる
	switch {
	case errors.Is(err, domain.ErrInvalidBooking):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrScheduleConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		slog.ErrorContext(
			r.Context(),
			"unexpected error",
			"error",
			err,
		)
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
	}
}

// #endregion
