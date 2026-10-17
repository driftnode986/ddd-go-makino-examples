package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/booking/domain"
)

// newTestBookingScheduledAt は scheduled に予約日時を置いた、確定前の予約を作る。
func newTestBookingScheduledAt(
	t *testing.T,
	scheduled time.Time,
) *domain.Booking {
	t.Helper()
	return newTestBookingAt(t, "10:00", "11:00", scheduled)
}

// #region expiration-test
func TestExpirationChecker_IsExpired(t *testing.T) {
	scheduled := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	booking := newTestBookingScheduledAt(t, scheduled)

	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"before scheduled", scheduled.Add(-1 * time.Hour), false},
		{"just at scheduled", scheduled, false},
		{
			"within grace period",
			scheduled.Add(20 * time.Minute),
			false,
		},
		{"past grace period", scheduled.Add(31 * time.Minute), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := domain.NewExpirationChecker(
				domain.FixedClock{T: tt.now},
			)
			if got := checker.IsExpired(booking); got != tt.want {
				t.Errorf("IsExpired: want=%v got=%v", tt.want, got)
			}
		})
	}
}

// #endregion

// 境界のケース (本に載せる表の残り)。30 分ちょうどはまだ失効しない。
func TestExpirationChecker_IsExpired_Boundaries(t *testing.T) {
	scheduled := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	booking := newTestBookingScheduledAt(t, scheduled)
	for _, tt := range []struct {
		name string
		now  time.Time
		want bool
	}{
		{
			"at grace deadline",
			scheduled.Add(30 * time.Minute),
			false,
		},
		{"far past", scheduled.Add(24 * time.Hour), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checker := domain.NewExpirationChecker(
				domain.FixedClock{T: tt.now},
			)
			if got := checker.IsExpired(booking); got != tt.want {
				t.Errorf("IsExpired: want=%v got=%v", tt.want, got)
			}
		})
	}
}
