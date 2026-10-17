package http_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/presentation/http"
)

// panic したリクエストもログに残ることを確かめる。
// Recoverer がロギングより外側にあると、panic がロギングを
// 素通りしてログの行が出ない。
func TestNewRouter_LogsPanickedRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	handler := httpx.NewBookingHandler(nil, nil)
	router := httpx.NewRouter(handler, logger)
	router.Get("/panic", func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if !strings.Contains(buf.String(), "path=/panic") {
		t.Errorf("panic したリクエストのログが無い: %q", buf.String())
	}
}
