package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/infrastructure/memory"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
	"github.com/google/uuid"
)

func TestPlaceBookingUseCase_SavesBooking(t *testing.T) {
	repo := memory.NewBookingRepository()
	memberID := uuid.NewString()
	stub := &stubLessonQuery{view: dto.LessonView{Name: "Yoga"}}
	now := func() time.Time {
		return time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC)
	}
	uc := application.NewPlaceBookingUseCase(repo, stub, now)

	out, err := uc.Execute(
		context.Background(),
		application.PlaceBookingInput{
			MemberID: memberID,
			LessonID: uuid.NewString(),
		},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if out.BookingID == "" {
		t.Fatal("BookingID is empty")
	}

	views, err := application.NewListBookingsQuery(repo).
		ByMember(context.Background(), memberID)
	if err != nil {
		t.Fatalf("ByMember() error = %v", err)
	}
	if len(views) != 1 || views[0].ID != out.BookingID {
		t.Errorf("views = %+v", views)
	}
}

func TestPlaceBookingUseCase_LessonNotAvailable(t *testing.T) {
	stub := &stubLessonQuery{err: errors.New("lesson not found")}
	uc := application.NewPlaceBookingUseCase(
		memory.NewBookingRepository(),
		stub,
		time.Now,
	)

	_, err := uc.Execute(
		context.Background(),
		application.PlaceBookingInput{
			MemberID: uuid.NewString(),
			LessonID: uuid.NewString(),
		},
	)
	if err == nil {
		t.Fatal("Execute() expected error, got nil")
	}
}

func TestPlaceBookingUseCase_InvalidMemberID(t *testing.T) {
	uc := application.NewPlaceBookingUseCase(
		memory.NewBookingRepository(),
		&stubLessonQuery{},
		time.Now,
	)
	_, err := uc.Execute(
		context.Background(),
		application.PlaceBookingInput{
			MemberID: "not-uuid",
			LessonID: uuid.NewString(),
		},
	)
	if err == nil {
		t.Fatal("Execute() expected error, got nil")
	}
}
