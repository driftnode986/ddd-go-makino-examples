package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch04-entity/internal/booking/domain"
)

func TestNewID_RejectsTooShort(t *testing.T) {
	_, err := domain.NewID[domain.BookingTag]("short")
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("want ErrInvalidID, got %v", err)
	}
}

func TestID_StringEqualsIsZero(t *testing.T) {
	a, _ := domain.NewID[domain.BookingTag]("book-0001")
	b, _ := domain.NewID[domain.BookingTag]("book-0001")
	if a.String() != "book-0001" {
		t.Fatalf("String(): got %q", a.String())
	}
	if !a.Equals(b) {
		t.Fatal("same value should be equal")
	}
	var zero domain.BookingID
	if !zero.IsZero() || a.IsZero() {
		t.Fatal("IsZero: want true only for the zero value")
	}
}

// ID がゼロ値の Booking は作れない
func TestNewBooking_RejectsZeroID(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	mid, _ := domain.NewID[domain.MemberTag]("memb-0001")
	lid, _ := domain.NewID[domain.LessonTag]("less-0001")
	_, err := domain.NewBooking(
		domain.BookingID{},
		mid,
		lid,
		now.Add(time.Hour),
		clock,
	)
	if err == nil {
		t.Fatal("want error for zero BookingID")
	}
}
