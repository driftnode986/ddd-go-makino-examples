package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

func newBooking(t *testing.T, at time.Time) *domain.Booking {
	t.Helper()
	b, err := domain.NewBooking("member-1", "lesson-1", at)
	if err != nil {
		t.Fatalf("NewBooking: %v", err)
	}
	return b
}

func TestNewBooking_RequiresFields(t *testing.T) {
	_, err := domain.NewBooking("", "lesson-1", time.Now().Add(48*time.Hour))
	if !errors.Is(err, domain.ErrInvalidBooking) {
		t.Fatalf("want ErrInvalidBooking, got %v", err)
	}
}

func TestNewBooking_StartsPendingWithFields(t *testing.T) {
	at := time.Now().Add(48 * time.Hour)
	b := newBooking(t, at)
	if b.Status() != domain.StatusPending {
		t.Errorf("Status: got %s", b.Status())
	}
	if b.MemberID() != "member-1" || b.LessonID() != "lesson-1" || !b.ScheduledAt().Equal(at) {
		t.Errorf("fields not set: %s %s %v", b.MemberID(), b.LessonID(), b.ScheduledAt())
	}
	if b.ID().String() == "" {
		t.Error("ID should be generated")
	}
}

func TestBooking_Confirm(t *testing.T) {
	b := newBooking(t, time.Now().Add(48*time.Hour))
	if err := b.Confirm(); err != nil {
		t.Fatalf("first Confirm: %v", err)
	}
	if b.Status() != domain.StatusConfirmed {
		t.Fatalf("Status: got %s", b.Status())
	}
	if err := b.Confirm(); !errors.Is(err, domain.ErrBookingAlreadyConfirmed) {
		t.Fatalf("second Confirm: want ErrBookingAlreadyConfirmed, got %v", err)
	}
}

func TestBooking_Cancel(t *testing.T) {
	t.Run("more than 24 hours before", func(t *testing.T) {
		b := newBooking(t, time.Now().Add(48*time.Hour))
		if err := b.Cancel("schedule changed"); err != nil {
			t.Fatalf("Cancel: %v", err)
		}
		if b.Status() != domain.StatusCancelled || b.CancelReason() != "schedule changed" {
			t.Fatalf("got status=%s reason=%q", b.Status(), b.CancelReason())
		}
	})
	t.Run("within 24 hours", func(t *testing.T) {
		b := newBooking(t, time.Now().Add(1*time.Hour))
		if err := b.Cancel("late"); !errors.Is(err, domain.ErrBookingCannotCancel) {
			t.Fatalf("want ErrBookingCannotCancel, got %v", err)
		}
	})
	t.Run("completed", func(t *testing.T) {
		b := newBooking(t, time.Now().Add(48*time.Hour))
		if err := b.Confirm(); err != nil {
			t.Fatal(err)
		}
		if err := b.MarkCompleted(); err != nil {
			t.Fatal(err)
		}
		if err := b.Cancel("too late"); !errors.Is(err, domain.ErrBookingCannotCancel) {
			t.Fatalf("want ErrBookingCannotCancel, got %v", err)
		}
	})
}

func TestBooking_MarkCompleted_OnlyFromConfirmed(t *testing.T) {
	b := newBooking(t, time.Now().Add(48*time.Hour))
	if err := b.MarkCompleted(); !errors.Is(err, domain.ErrBookingCannotComplete) {
		t.Fatalf("pending: want ErrBookingCannotComplete, got %v", err)
	}
}

func TestBooking_Reschedule(t *testing.T) {
	b := newBooking(t, time.Now().Add(48*time.Hour))
	newAt := time.Now().Add(72 * time.Hour)
	if err := b.Reschedule(newAt); err != nil {
		t.Fatalf("Reschedule: %v", err)
	}
	if !b.ScheduledAt().Equal(newAt) {
		t.Fatalf("ScheduledAt: got %v", b.ScheduledAt())
	}
	if err := b.Cancel("x"); err != nil {
		t.Fatal(err)
	}
	if err := b.Reschedule(newAt); !errors.Is(err, domain.ErrBookingCannotReschedule) {
		t.Fatalf("cancelled: want ErrBookingCannotReschedule, got %v", err)
	}
}

func TestScheduleConflictError_IsAndAs(t *testing.T) {
	var err error = &domain.ScheduleConflictError{
		ExistingBookingID: newBooking(t, time.Now().Add(48*time.Hour)).ID(),
		RequestedAt:       time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
	}
	if !errors.Is(err, domain.ErrScheduleConflict) {
		t.Fatal("errors.Is(err, ErrScheduleConflict) should be true")
	}
	if errors.Is(err, domain.ErrBookingNotFound) {
		t.Fatal("should not match another sentinel")
	}
	var conflict *domain.ScheduleConflictError
	if !errors.As(err, &conflict) {
		t.Fatal("errors.As should find *ScheduleConflictError")
	}
	if conflict.Error() == "" {
		t.Fatal("Error() should describe the conflict")
	}
}

func TestAreSlotsConflicting(t *testing.T) {
	at := func(h int) time.Time { return time.Date(2026, 6, 1, h, 0, 0, 0, time.UTC) }
	tests := []struct {
		name string
		a, b domain.TimeSlot
		want bool
	}{
		{"overlap", domain.TimeSlot{Start: at(10), End: at(12)}, domain.TimeSlot{Start: at(11), End: at(13)}, true},
		{"adjacent", domain.TimeSlot{Start: at(10), End: at(11)}, domain.TimeSlot{Start: at(11), End: at(12)}, false},
		{"separate", domain.TimeSlot{Start: at(9), End: at(10)}, domain.TimeSlot{Start: at(14), End: at(15)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domain.AreSlotsConflicting(tt.a, tt.b); got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestIDs_String(t *testing.T) {
	if domain.MemberID("m-1").String() != "m-1" || domain.LessonID("l-1").String() != "l-1" {
		t.Fatal("String() should return the raw value")
	}
}
