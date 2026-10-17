// #region request-id-middleware-a
package presentation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// #endregion
// #region request-id-middleware-b
// RequestIDMiddleware は X-Request-ID ヘッダを尊重しつつ、
// なければ新しい ID を生成して context に乗せる。
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = newRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// #endregion

// #region request-id-helpers
func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

func newRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// #endregion

// #region logging-middleware
// LoggingMiddleware は HTTP リクエスト/レスポンスを
// 構造化ログで記録する。
func LoggingMiddleware(
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				sw := &statusWriter{
					ResponseWriter: w,
					status:         http.StatusOK,
				}

				next.ServeHTTP(sw, r)

				id, _ := RequestIDFromContext(r.Context())
				logger.InfoContext(r.Context(), "http request",
					slog.String("request_id", id),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", sw.status),
					slog.Duration("elapsed", time.Since(start)),
				)
			},
		)
	}
}

// #endregion

// #region status-writer
// statusWriter は WriteHeader で渡されたステータスを保持する。
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

// #endregion

// #region recovery-middleware-a
// RecoveryMiddleware は panic を 500 に変換し、
// スタックトレースをログに記録する。
func RecoveryMiddleware(
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// #endregion
				// #region recovery-middleware-b
				defer func() {
					if rec := recover(); rec != nil {
						logger.ErrorContext(
							r.Context(),
							"panic recovered",
							slog.Any("panic", rec),
							slog.String(
								"stack",
								string(debug.Stack()),
							),
						)
						writeError(
							w,
							http.StatusInternalServerError,
							"INTERNAL_ERROR",
							"internal server error",
						)
					}
				}()
				next.ServeHTTP(w, r)
			},
		)
	}
}

// #endregion
