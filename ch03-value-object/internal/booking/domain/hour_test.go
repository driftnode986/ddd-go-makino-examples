package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/domain"
)

func TestNewHour_RejectsOutOfBusinessHours(t *testing.T) {
	tests := []struct {
		name    string
		input   time.Time
		wantErr error
	}{
		{
			"too early",
			time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC),
			domain.ErrOutOfBusinessHours,
		},
		{
			"too late",
			time.Date(2026, 6, 1, 22, 0, 0, 0, time.UTC),
			domain.ErrOutOfBusinessHours,
		},
		{
			"has minutes",
			time.Date(2026, 6, 1, 10, 30, 0, 0, time.UTC),
			domain.ErrHourHasMinutes,
		},
		{
			"valid 9",
			time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC),
			nil,
		},
		{
			"valid 21",
			time.Date(2026, 6, 1, 21, 0, 0, 0, time.UTC),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewHour(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewHour: want=%v got=%v", tt.wantErr, err)
			}
		})
	}
}

func TestHour_IsInThePast(t *testing.T) {
	frozen := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return frozen }

	past, _ := domain.NewHour(
		time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC),
	)
	if !past.IsInThePast(clock) {
		t.Errorf("9:00 should be in the past at 12:00")
	}

	future, _ := domain.NewHour(
		time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC),
	)
	if future.IsInThePast(clock) {
		t.Errorf("15:00 should not be in the past at 12:00")
	}
}
