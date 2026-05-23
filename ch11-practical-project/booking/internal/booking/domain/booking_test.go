// internal/booking/domain/booking_test.go
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
