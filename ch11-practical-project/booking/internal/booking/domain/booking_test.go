// #region confirm-from-pending
package domain_test

import (
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/testdata"
)

func TestBooking_Confirm_FromPending(t *testing.T) {
	b := testdata.NewBookingBuilder().Build()

	if err := b.Confirm(); err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	if b.Status() != domain.BookingStatusConfirmed {
		t.Errorf("status = %v, want %v", b.Status(), domain.BookingStatusConfirmed)
	}
	if b.Version() != 2 {
		t.Errorf("version = %d, want 2", b.Version())
	}
}

// #endregion

func TestBooking_Confirm_AlreadyCancelled(t *testing.T) {
	b := testdata.NewBookingBuilder().
		WithStatus(domain.BookingStatusCancelled).
		Build()

	if err := b.Confirm(); err == nil {
		t.Fatal("Confirm() expected error, got nil")
	}
}

func TestBooking_Cancel_FromPending(t *testing.T) {
	b := testdata.NewBookingBuilder().Build()

	if err := b.Cancel(); err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if b.Status() != domain.BookingStatusCancelled {
		t.Errorf("status = %v, want %v", b.Status(), domain.BookingStatusCancelled)
	}
}

func TestBooking_Cancel_FromConfirmed(t *testing.T) {
	// #region builder-usage
	booking := testdata.NewBookingBuilder().
		WithStatus(domain.BookingStatusConfirmed).
		Build()
	// #endregion

	if err := booking.Cancel(); err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if booking.Status() != domain.BookingStatusCancelled {
		t.Errorf("status = %v, want %v", booking.Status(), domain.BookingStatusCancelled)
	}
	if booking.Version() != 2 {
		t.Errorf("version = %d, want 2", booking.Version())
	}
	if err := booking.Cancel(); err == nil {
		t.Error("second Cancel() expected error, got nil")
	}
}

func TestBookingBuilder_WithMember(t *testing.T) {
	member := domain.NewMemberID()
	b := testdata.NewBookingBuilder().WithMember(member).Build()
	if b.MemberID() != member {
		t.Errorf("member = %v, want %v", b.MemberID(), member)
	}
}

func TestBookingIDFrom(t *testing.T) {
	id := domain.NewBookingID()
	got, err := domain.BookingIDFrom(id.String())
	if err != nil || got != id {
		t.Errorf("BookingIDFrom(%q) = %v, %v", id, got, err)
	}
	if _, err := domain.BookingIDFrom("not-uuid"); err == nil {
		t.Error("BookingIDFrom(not-uuid) expected error, got nil")
	}
}
