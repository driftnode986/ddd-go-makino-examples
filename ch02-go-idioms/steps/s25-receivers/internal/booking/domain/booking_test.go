package domain

import (
	"errors"
	"testing"
)

func TestBooking_Confirm_OnlyFromPending(t *testing.T) {
	b := &Booking{status: StatusPending}
	if err := b.Confirm(); err != nil {
		t.Fatal(err)
	}
	if b.Status() != StatusConfirmed {
		t.Fatalf("Status: got %s", b.Status())
	}
	if err := b.Confirm(); !errors.Is(err, ErrBookingAlreadyConfirmed) {
		t.Fatalf("want ErrBookingAlreadyConfirmed, got %v", err)
	}
	var zero BookingID
	if b.ID() != zero {
		t.Fatal("ID: want the zero value")
	}
}
