package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
)

func newBooking(t *testing.T) *domain.Booking {
	t.Helper()
	id, _ := domain.NewBookingID("bk_001")
	mid, _ := domain.NewMemberID("mem_001")
	lid, _ := domain.NewLessonID("lsn_001")
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	return domain.NewBooking(id, mid, lid, at)
}

func TestBooking_NewBooking_DefaultsToPendingVersion1(t *testing.T) {
	b := newBooking(t)
	if b.Status() != domain.StatusPending {
		t.Errorf(
			"status = %v, want %v",
			b.Status(),
			domain.StatusPending,
		)
	}
	if b.Version() != 1 {
		t.Errorf("version = %d, want 1", b.Version())
	}
}

func TestBooking_Confirm_IncrementsVersion(t *testing.T) {
	b := newBooking(t)
	if err := b.Confirm(); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if b.Status() != domain.StatusConfirmed {
		t.Errorf(
			"status = %v, want %v",
			b.Status(),
			domain.StatusConfirmed,
		)
	}
	if b.Version() != 2 {
		t.Errorf("version = %d, want 2", b.Version())
	}
}

func TestBooking_Cancel_FromConfirmed(t *testing.T) {
	b := newBooking(t)
	_ = b.Confirm()
	if err := b.Cancel(); err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if b.Status() != domain.StatusCancelled {
		t.Errorf(
			"status = %v, want %v",
			b.Status(),
			domain.StatusCancelled,
		)
	}
}

func TestBooking_Cancel_AlreadyCancelled(t *testing.T) {
	b := newBooking(t)
	_ = b.Cancel()
	err := b.Cancel()
	if !errors.Is(err, domain.ErrCannotCancel) {
		t.Errorf("got %v, want ErrCannotCancel", err)
	}
}

func TestBooking_Confirm_FromCancelled_Fails(t *testing.T) {
	b := newBooking(t)
	_ = b.Cancel()
	err := b.Confirm()
	if !errors.Is(err, domain.ErrCannotCancel) {
		t.Errorf("got %v, want ErrCannotCancel", err)
	}
}
