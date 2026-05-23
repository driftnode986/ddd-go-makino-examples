// internal/booking/presentation/http/booking_handler_test.go
package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/infrastructure/memory"
	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/presentation/http"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
	"github.com/google/uuid"
)

type stubLessonQuery struct{ view dto.LessonView }

func (s *stubLessonQuery) FindLessonView(_ context.Context,
	_ string) (dto.LessonView, error) {
	return s.view, nil
}

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	repo := memory.NewBookingRepository()
	stub := &stubLessonQuery{view: dto.LessonView{ID: uuid.NewString(), Name: "Yoga"}}
	now := func() time.Time { return time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC) }
	placeUC := application.NewPlaceBookingUseCase(repo, stub, now)
	listQ := application.NewListBookingsQuery(repo)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := httpx.NewBookingHandler(placeUC, listQ)
	router := httpx.NewRouter(handler, logger)
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return srv
}

func TestPlaceBooking_HTTP(t *testing.T) {
	srv := setupTestServer(t)
	body := strings.NewReader(`{
		"member_id": "` + uuid.NewString() + `",
		"lesson_id": "` + uuid.NewString() + `"
	}`)

	res, err := http.Post(srv.URL+"/booking/bookings", "application/json", body)
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(res.Body)
		t.Fatalf("status = %d, body = %s", res.StatusCode, raw)
	}

	var out struct {
		BookingID string `json:"BookingID"`
		Status    string `json:"Status"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Status != "pending" {
		t.Errorf("status = %q, want pending", out.Status)
	}
}

func TestPlaceBooking_HTTP_InvalidInput(t *testing.T) {
	srv := setupTestServer(t)
	body := bytes.NewBufferString(`{"member_id":"not-uuid","lesson_id":"not-uuid"}`)

	res, err := http.Post(srv.URL+"/booking/bookings", "application/json", body)
	if err != nil {
		t.Fatalf("POST error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}
