package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/domain"
)

func TestEmail_Domain_IsLowercased(t *testing.T) {
	e, err := domain.NewEmail("user@Example.COM")
	if err != nil {
		t.Fatal(err)
	}
	if got := e.Domain(); got != "example.com" {
		t.Fatalf("Domain(): got %q want %q", got, "example.com")
	}
}

func TestHour_Time_ReturnsGivenTime(t *testing.T) {
	at := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	h, err := domain.NewHour(at)
	if err != nil {
		t.Fatal(err)
	}
	if !h.Time().Equal(at) {
		t.Fatalf("Time(): got %v want %v", h.Time(), at)
	}
}

func TestMoney_AccessorsAndString(t *testing.T) {
	m, err := domain.NewMoney(1500, domain.JPY)
	if err != nil {
		t.Fatal(err)
	}
	if m.Amount() != 1500 || m.Currency() != domain.JPY {
		t.Fatalf("got %d %s", m.Amount(), m.Currency())
	}
	if got := m.String(); got != "1500 JPY" {
		t.Fatalf("String(): got %q", got)
	}
}

func TestNewMoney_RequiresCurrency(t *testing.T) {
	_, err := domain.NewMoney(100, "")
	if !errors.Is(err, domain.ErrCurrencyRequired) {
		t.Fatalf("want ErrCurrencyRequired, got %v", err)
	}
}

func TestID_String_ReturnsValue(t *testing.T) {
	id, err := domain.NewID[domain.Booking]("booking-0001")
	if err != nil {
		t.Fatal(err)
	}
	if got := id.String(); got != "booking-0001" {
		t.Fatalf("String(): got %q", got)
	}
}
