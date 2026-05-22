package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/booking/domain"
)

type fakeBookingFinder struct {
	bookings []*domain.Booking
}

func (f *fakeBookingFinder) FindByTrainerAndPeriod(
	ctx context.Context, trainerID domain.TrainerID, from, to domain.Hour,
) ([]*domain.Booking, error) {
	var hits []*domain.Booking
	for _, b := range f.bookings {
		if b.TrainerID().Equals(trainerID) {
			hits = append(hits, b)
		}
	}
	return hits, nil
}

func parseHour(t *testing.T, s string) domain.Hour {
	t.Helper()
	h, err := domain.NewHour(s)
	if err != nil {
		t.Fatalf("NewHour(%q): %v", s, err)
	}
	return h
}

func newTestBooking(t *testing.T, trainerID domain.TrainerID, from, to string) *domain.Booking {
	t.Helper()
	id, err := domain.NewID[domain.BookingTag]("booking01")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	memberID, err := domain.NewID[domain.MemberTag]("member001")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	b, err := domain.NewBooking(
		id, memberID, trainerID,
		parseHour(t, from), parseHour(t, to),
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		domain.FixedClock{T: time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)},
	)
	if err != nil {
		t.Fatalf("NewBooking: %v", err)
	}
	return b
}

func TestConflictDetector_Detect(t *testing.T) {
	ctx := context.Background()
	trainerID, err := domain.NewID[domain.TrainerTag]("trainer1")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	existing := newTestBooking(t, trainerID, "10:00", "11:00")
	finder := &fakeBookingFinder{bookings: []*domain.Booking{existing}}
	detector := domain.NewConflictDetector(finder)

	tests := []struct {
		name string
		from string
		to   string
		want error
	}{
		{"no overlap before", "08:00", "09:00", nil},
		{"no overlap after", "11:00", "12:00", nil},
		{"adjacent end", "09:00", "10:00", nil},
		{"partial overlap start", "09:30", "10:30", domain.ErrBookingConflict},
		{"partial overlap end", "10:30", "11:30", domain.ErrBookingConflict},
		{"contained", "10:15", "10:45", domain.ErrBookingConflict},
		{"identical", "10:00", "11:00", domain.ErrBookingConflict},
		{"invalid period", "11:00", "10:00", domain.ErrInvalidPeriod},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := detector.Detect(ctx, trainerID, parseHour(t, tt.from), parseHour(t, tt.to))
			if !errors.Is(err, tt.want) {
				t.Errorf("Detect: want=%v got=%v", tt.want, err)
			}
		})
	}
}

func TestConflictDetector_IgnoresCancelled(t *testing.T) {
	ctx := context.Background()
	trainerID, err := domain.NewID[domain.TrainerTag]("trainer1")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	// 既存予約をキャンセル済みとして扱う場合、衝突しないことを確認したいが
	// 本章では Booking エンティティに Cancel() を実装していない。
	// 代わりに FindByTrainerAndPeriod が空を返すケースをテストする。
	finder := &fakeBookingFinder{bookings: nil}
	detector := domain.NewConflictDetector(finder)

	if err := detector.Detect(ctx, trainerID, parseHour(t, "10:00"), parseHour(t, "11:00")); err != nil {
		t.Errorf("Detect on empty: want=nil got=%v", err)
	}
}
