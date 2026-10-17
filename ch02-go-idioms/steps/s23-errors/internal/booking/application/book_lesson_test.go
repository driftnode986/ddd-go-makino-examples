package application_test

import (
	"context"
	"errors"
	"iter"
	"strings"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/steps/s23-errors/internal/booking/application"
)

type saveOnlyRepo struct{ err error }

func (r saveOnlyRepo) Save(context.Context, *domain.Booking) error { return r.err }
func (saveOnlyRepo) FindByID(context.Context, domain.BookingID) (*domain.Booking, error) {
	return nil, domain.ErrBookingNotFound
}
func (saveOnlyRepo) AllByMember(context.Context, domain.MemberID) iter.Seq2[*domain.Booking, error] {
	return func(func(*domain.Booking, error) bool) {}
}

func input() application.BookLessonInput {
	return application.BookLessonInput{MemberID: "member-1", LessonID: "lesson-1", ScheduledAt: time.Now().Add(48 * time.Hour)}
}

func TestExecute_WrapsConflictWithExistingID(t *testing.T) {
	existing, err := domain.NewBooking("member-2", "lesson-1", time.Now().Add(48*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	repo := saveOnlyRepo{err: &domain.ScheduleConflictError{ExistingBookingID: existing.ID()}}
	err = application.NewBookLessonUseCase(repo).Execute(context.Background(), input())
	if !errors.Is(err, domain.ErrScheduleConflict) {
		t.Fatalf("want ErrScheduleConflict, got %v", err)
	}
	if !strings.HasPrefix(err.Error(), "conflict with "+existing.ID().String()) {
		t.Fatalf("message should name the existing booking: %v", err)
	}
}

func TestExecute_WrapsOtherSaveErrors(t *testing.T) {
	boom := errors.New("db down")
	err := application.NewBookLessonUseCase(saveOnlyRepo{err: boom}).Execute(context.Background(), input())
	if !errors.Is(err, boom) || !strings.HasPrefix(err.Error(), "save booking: ") {
		t.Fatalf("got %v", err)
	}
}

func TestExecute_Succeeds(t *testing.T) {
	if err := application.NewBookLessonUseCase(saveOnlyRepo{}).Execute(context.Background(), input()); err != nil {
		t.Fatal(err)
	}
}

func TestExecute_RejectsInvalidInput(t *testing.T) {
	in := input()
	in.LessonID = ""
	err := application.NewBookLessonUseCase(saveOnlyRepo{}).Execute(context.Background(), in)
	if !errors.Is(err, domain.ErrInvalidBooking) {
		t.Fatalf("got %v", err)
	}
}
