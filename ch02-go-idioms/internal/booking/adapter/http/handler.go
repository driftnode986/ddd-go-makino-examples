// #region handler-imports
package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
	identity "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/identity/domain"
)

// #endregion

// #region handle-create-booking
// アプリケーション層でドメインエラーを HTTP ステータスに対応させる
func (h *Handler) handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Execute(r.Context(), buildInput(r))
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}

	switch {
	case errors.Is(err, identity.ErrInvalidEmail),
		errors.Is(err, identity.ErrInvalidDisplayName):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrBookingNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrScheduleConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		slog.ErrorContext(r.Context(), "unexpected error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// #endregion
