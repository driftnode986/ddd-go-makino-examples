// #region error-response-a
package presentation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// #endregion
// #region error-response-b
func writeError(
	w http.ResponseWriter,
	status int,
	code, message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).
		Encode(errorResponse{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// #endregion

// #region write-application-error
func writeApplicationError(
	w http.ResponseWriter,
	logger *slog.Logger,
	r *http.Request,
	err error,
) {
	var appErr *application.ApplicationError
	if errors.As(err, &appErr) {
		status := mapApplicationCodeToHTTP(appErr.Code)
		if status >= 500 {
			logger.ErrorContext(r.Context(), "application error",
				slog.String("code", appErr.Code),
				slog.Any("err", err),
			)
		}
		writeError(w, status, appErr.Code, appErr.Message)
		return
	}

	// ApplicationError でラップされていない場合は内部エラー扱い
	logger.ErrorContext(r.Context(), "unexpected error",
		slog.Any("err", err),
	)
	writeError(w, http.StatusInternalServerError,
		"INTERNAL_ERROR", "internal server error")
}

// #endregion

// #region map-code
func mapApplicationCodeToHTTP(code string) int {
	switch code {
	case "INVALID_INPUT":
		return http.StatusBadRequest
	case "NOT_FOUND":
		return http.StatusNotFound
	case "CONFLICT", "INSUFFICIENT_STOCK":
		return http.StatusConflict
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// #endregion
