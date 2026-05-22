package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/booking/domain"
)

func TestExpirationChecker_IsExpired(t *testing.T) {
	scheduled := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	trainerID, err := domain.NewID[domain.TrainerTag]("trainer1")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	id, err := domain.NewID[domain.BookingTag]("booking01")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	memberID, err := domain.NewID[domain.MemberTag]("member001")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}

	booking, err := domain.NewBooking(
		id, memberID, trainerID,
		parseHour(t, "10:00"), parseHour(t, "11:00"),
		scheduled,
		domain.FixedClock{T: time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)},
	)
	if err != nil {
		t.Fatalf("NewBooking: %v", err)
	}

	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"before scheduled", scheduled.Add(-1 * time.Hour), false},
		{"just at scheduled", scheduled, false},
		{"within grace period", scheduled.Add(20 * time.Minute), false},
		{"at grace deadline", scheduled.Add(30 * time.Minute), false},
		{"past grace period", scheduled.Add(31 * time.Minute), true},
		{"far past", scheduled.Add(24 * time.Hour), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := domain.NewExpirationChecker(domain.FixedClock{T: tt.now})
			if got := checker.IsExpired(booking); got != tt.want {
				t.Errorf("IsExpired: want=%v got=%v", tt.want, got)
			}
		})
	}
}
