package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
	identity "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/identity/domain"
)

type stubUseCase struct {
	err error
	got application.BookLessonInput
}

func (s *stubUseCase) Execute(
	_ context.Context,
	in application.BookLessonInput,
) error {
	s.got = in
	return s.err
}

const body = `{"member_id":"member-1","lesson_id":"lesson-1",` +
	`"scheduled_at":"2026-06-01T10:00:00Z"}`

func TestHandleCreateBooking_MapsErrorsToStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"created", nil, http.StatusCreated},
		{
			"invalid email",
			identity.ErrInvalidEmail,
			http.StatusBadRequest,
		},
		{
			"invalid display name",
			identity.ErrInvalidDisplayName,
			http.StatusBadRequest,
		},
		{
			"not found",
			domain.ErrBookingNotFound,
			http.StatusNotFound,
		},
		{
			"conflict (wrapped error type)",
			&domain.ScheduleConflictError{},
			http.StatusConflict,
		},
		{
			"unexpected",
			errors.New("db down"),
			http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &stubUseCase{err: tt.err}
			rec := httptest.NewRecorder()
			NewHandler(uc).ServeHTTP(
				rec,
				httptest.NewRequest(
					http.MethodPost,
					"/bookings",
					strings.NewReader(body),
				),
			)
			if rec.Code != tt.want {
				t.Fatalf("status: got %d want %d", rec.Code, tt.want)
			}
		})
	}
}

func TestBuildInput_DecodesRequest(t *testing.T) {
	uc := &stubUseCase{}
	NewHandler(uc).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(
			http.MethodPost,
			"/bookings",
			strings.NewReader(body),
		),
	)
	if uc.got.MemberID != "member-1" ||
		uc.got.LessonID != "lesson-1" ||
		uc.got.ScheduledAt.Hour() != 10 {
		t.Fatalf("input not decoded: %+v", uc.got)
	}
}

func TestServeHTTP_RejectsOtherMethods(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHandler(&stubUseCase{}).ServeHTTP(
		rec,
		httptest.NewRequest(http.MethodGet, "/bookings", nil),
	)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: got %d", rec.Code)
	}
}

// 本物のユースケースとつないでも型が合う
var _ bookLessonExecutor = (*application.BookLessonUseCase)(nil)
