// #region stub-lesson-query
package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/infrastructure/memory"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
	"github.com/google/uuid"
)

type stubLessonQuery struct {
	view dto.LessonView
	err  error
}

func (s *stubLessonQuery) FindLessonView(_ context.Context,
	_ string) (dto.LessonView, error) {
	return s.view, s.err
}

// #endregion

// #region place-booking-success
func TestPlaceBookingUseCase_Success(t *testing.T) {
	repo := memory.NewBookingRepository()
	lessonID := uuid.NewString()
	memberID := uuid.NewString()
	stub := &stubLessonQuery{view: dto.LessonView{ID: lessonID, Name: "Yoga"}}
	now := func() time.Time { return time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC) }

	uc := application.NewPlaceBookingUseCase(repo, stub, now)

	out, err := uc.Execute(context.Background(), application.PlaceBookingInput{
		MemberID: memberID,
		LessonID: lessonID,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if out.Status != "pending" {
		t.Errorf("status = %q, want %q", out.Status, "pending")
	}
}

// #endregion
