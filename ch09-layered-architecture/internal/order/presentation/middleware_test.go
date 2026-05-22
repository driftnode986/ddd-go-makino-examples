package presentation_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/presentation"
)

func TestRequestIDMiddleware_GeneratesNewID(t *testing.T) {
	called := false
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if id, ok := presentation.RequestIDFromContext(r.Context()); ok {
			capturedID = id
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := presentation.RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler was not called")
	}
	if capturedID == "" {
		t.Fatal("request id was not propagated to context")
	}
	respID := rec.Header().Get("X-Request-ID")
	if respID == "" {
		t.Fatal("X-Request-ID response header is empty")
	}
	if respID != capturedID {
		t.Errorf("header id = %q, context id = %q", respID, capturedID)
	}
}

func TestRequestIDMiddleware_PreservesExistingID(t *testing.T) {
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, ok := presentation.RequestIDFromContext(r.Context()); ok {
			capturedID = id
		}
	})

	handler := presentation.RequestIDMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req-12345")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if capturedID != "req-12345" {
		t.Errorf("captured = %q, want %q", capturedID, "req-12345")
	}
	if rec.Header().Get("X-Request-ID") != "req-12345" {
		t.Errorf("response header = %q, want %q",
			rec.Header().Get("X-Request-ID"), "req-12345")
	}
}

func TestRecoveryMiddleware_RecoversFromPanic(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic for test")
	})

	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	handler := presentation.RecoveryMiddleware(logger)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic should have been recovered: %v", r)
		}
	}()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
