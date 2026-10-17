package after_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/after"
)

func newBooking(t *testing.T) *after.Booking {
	t.Helper()
	b, err := after.NewBooking("b-1", "m-1")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestNewBooking_ReturnsErrorOnEmptyID(t *testing.T) {
	if _, err := after.NewBooking("", "m-1"); !errors.Is(
		err,
		after.ErrBookingIDRequired,
	) {
		t.Fatalf("want ErrBookingIDRequired, got %v", err)
	}
}

func TestBooking_Cancel_RecordsReasonInOneCall(t *testing.T) {
	b := newBooking(t)
	if err := b.Cancel("schedule changed"); err != nil {
		t.Fatal(err)
	}
	if b.Status() != after.StatusCancelled ||
		b.CancelReason() != "schedule changed" {
		t.Fatalf(
			"got status=%s reason=%q",
			b.Status(),
			b.CancelReason(),
		)
	}
	if b.ID() != "b-1" || b.MemberID() != "m-1" {
		t.Fatalf("got id=%s member=%s", b.ID(), b.MemberID())
	}
}

func TestBooking_RulesStopInvalidTransitions(t *testing.T) {
	b := newBooking(t)
	if err := b.MarkCompleted(); !errors.Is(
		err,
		after.ErrCannotComplete,
	) {
		t.Fatalf("complete from pending: got %v", err)
	}
	if err := b.Confirm(); err != nil {
		t.Fatal(err)
	}
	if err := b.Confirm(); !errors.Is(err, after.ErrCannotConfirm) {
		t.Fatalf("confirm twice: got %v", err)
	}
	if err := b.MarkCompleted(); err != nil {
		t.Fatal(err)
	}
	if err := b.Cancel("too late"); !errors.Is(
		err,
		after.ErrCannotCancel,
	) {
		t.Fatalf("cancel after completed: got %v", err)
	}
}
