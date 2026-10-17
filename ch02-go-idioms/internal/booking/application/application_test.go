package application_test

import (
	"context"
	"errors"
	"iter"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// fakeRepo は BookingRepository のテスト用の実装。
type fakeRepo struct {
	saved   []*domain.Booking
	saveErr error
	all     []*domain.Booking
	iterErr error
}

func (r *fakeRepo) Save(_ context.Context, b *domain.Booking) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.saved = append(r.saved, b)
	return nil
}

func (r *fakeRepo) FindByID(_ context.Context, id domain.BookingID) (*domain.Booking, error) {
	for _, b := range r.saved {
		if b.ID() == id {
			return b, nil
		}
	}
	return nil, domain.ErrBookingNotFound
}

func (r *fakeRepo) AllByMember(_ context.Context, memberID domain.MemberID) iter.Seq2[*domain.Booking, error] {
	return func(yield func(*domain.Booking, error) bool) {
		for _, b := range r.all {
			if b.MemberID() != memberID {
				continue
			}
			if !yield(b, nil) {
				return
			}
		}
		if r.iterErr != nil {
			yield(nil, r.iterErr)
		}
	}
}

func input() application.BookLessonInput {
	return application.BookLessonInput{
		MemberID:    "member-1",
		LessonID:    "lesson-1",
		ScheduledAt: time.Now().Add(48 * time.Hour),
	}
}

func TestBookLessonUseCase_Execute_SavesBooking(t *testing.T) {
	repo := &fakeRepo{}
	if err := application.NewBookLessonUseCase(repo).Execute(context.Background(), input()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(repo.saved) != 1 || repo.saved[0].Status() != domain.StatusPending {
		t.Fatalf("want 1 pending booking saved, got %v", repo.saved)
	}
	if _, err := repo.FindByID(context.Background(), repo.saved[0].ID()); err != nil {
		t.Fatalf("FindByID: %v", err)
	}
}

func TestBookLessonUseCase_Execute_ReturnsSaveError(t *testing.T) {
	conflict := &domain.ScheduleConflictError{ExistingBookingID: mustBooking(t).ID()}
	repo := &fakeRepo{saveErr: conflict}
	err := application.NewBookLessonUseCase(repo).Execute(context.Background(), input())
	if !errors.Is(err, domain.ErrScheduleConflict) {
		t.Fatalf("want ErrScheduleConflict, got %v", err)
	}
}

func TestBookLessonUseCase_Execute_RejectsInvalidInput(t *testing.T) {
	in := input()
	in.MemberID = ""
	err := application.NewBookLessonUseCase(&fakeRepo{}).Execute(context.Background(), in)
	if !errors.Is(err, domain.ErrInvalidBooking) {
		t.Fatalf("want ErrInvalidBooking, got %v", err)
	}
}

func mustBooking(t *testing.T) *domain.Booking {
	t.Helper()
	b, err := domain.NewBooking("member-1", "lesson-1", time.Now().Add(48*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestProcessMemberBookings(t *testing.T) {
	a, b := mustBooking(t), mustBooking(t)
	other, err := domain.NewBooking("member-2", "lesson-1", time.Now().Add(48*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	repo := &fakeRepo{all: []*domain.Booking{a, other, b}}

	var got []*domain.Booking
	err = application.ProcessMemberBookings(context.Background(), repo, "member-1",
		func(bk *domain.Booking) error { got = append(got, bk); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("want the 2 bookings of member-1 in order, got %d", len(got))
	}

	t.Run("stops at the first process error", func(t *testing.T) {
		stop := errors.New("stop")
		calls := 0
		err := application.ProcessMemberBookings(context.Background(), repo, "member-1",
			func(*domain.Booking) error { calls++; return stop })
		if !errors.Is(err, stop) || calls != 1 {
			t.Fatalf("want stop after 1 call, got err=%v calls=%d", err, calls)
		}
	})
	t.Run("wraps the iteration error", func(t *testing.T) {
		boom := errors.New("db down")
		r := &fakeRepo{iterErr: boom}
		err := application.ProcessMemberBookings(context.Background(), r, "member-1",
			func(*domain.Booking) error { return nil })
		if !errors.Is(err, boom) {
			t.Fatalf("want wrapped iteration error, got %v", err)
		}
	})
}
