package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/booking/domain"
)

func TestFixedClock_Now_ReturnsT(t *testing.T) {
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	var c domain.Clock = domain.FixedClock{T: at}
	if !c.Now().Equal(at) {
		t.Fatalf("Now(): got %v want %v", c.Now(), at)
	}
}

func TestSystemClock_Now_ReturnsCurrentTime(t *testing.T) {
	var c domain.Clock = domain.SystemClock{}
	before := time.Now()
	got := c.Now()
	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("Now(): %v is not between %v and %v", got, before, after)
	}
}
