package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/booking/domain"
)

// #region fake-finder
type fakeBookingFinder struct {
	bookings []*domain.Booking
}

func (f *fakeBookingFinder) FindByTrainerAndPeriod(
	ctx context.Context,
	trainerID domain.TrainerID,
	from, to domain.Hour,
) ([]*domain.Booking, error) {
	var hits []*domain.Booking
	for _, b := range f.bookings {
		if b.TrainerID().Equals(trainerID) && overlaps(b, from, to) {
			hits = append(hits, b)
		}
	}
	return hits, nil
}

// #endregion

// overlaps は予約 b の時間帯が [from, to) と重なるかを返す (リポジトリの期間検索を模す)。
func overlaps(b *domain.Booking, from, to domain.Hour) bool {
	return domain.HoursOverlap(from, to, b.From(), b.To())
}

func parseHour(t *testing.T, s string) domain.Hour {
	t.Helper()
	h, err := domain.NewHour(s)
	if err != nil {
		t.Fatalf("NewHour(%q): %v", s, err)
	}
	return h
}

func newTestTrainerID(t *testing.T) domain.TrainerID {
	t.Helper()
	trainerID, err := domain.NewID[domain.TrainerTag]("trainer1")
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	return trainerID
}

// newTestBooking は newTestTrainerID のトレーナーの、from から to までの予約を作る。
func newTestBooking(t *testing.T, from, to string) *domain.Booking {
	t.Helper()
	return newTestBookingAt(
		t,
		from,
		to,
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	)
}

func newTestBookingAt(
	t *testing.T,
	from, to string,
	scheduledAt time.Time,
) *domain.Booking {
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
		id,
		memberID,
		newTestTrainerID(t),
		parseHour(t, from),
		parseHour(t, to),
		scheduledAt,
		domain.FixedClock{
			T: time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC),
		},
	)
	if err != nil {
		t.Fatalf("NewBooking: %v", err)
	}
	return b
}

// #region detect-test
func TestConflictDetector_Detect(t *testing.T) {
	ctx := context.Background()
	trainerID := newTestTrainerID(t)
	existing := newTestBooking(t, "10:00", "11:00")
	finder := &fakeBookingFinder{
		bookings: []*domain.Booking{existing},
	}
	detector := domain.NewConflictDetector(finder)

	tests := []struct {
		name string
		from string
		to   string
		want error
	}{
		{"no overlap before", "08:00", "09:00", nil},
		{"no overlap after", "11:00", "12:00", nil},
		{
			"partial overlap start",
			"09:30",
			"10:30",
			domain.ErrBookingConflict,
		},
		{"contained", "10:15", "10:45", domain.ErrBookingConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := detector.Detect(
				ctx,
				trainerID,
				parseHour(t, tt.from),
				parseHour(t, tt.to),
			)
			if !errors.Is(err, tt.want) {
				t.Errorf("Detect: want=%v got=%v", tt.want, err)
			}
		})
	}
}

// #endregion

// 境界のケース (本に載せる表の残り)。
func TestConflictDetector_Detect_Boundaries(t *testing.T) {
	ctx := context.Background()
	trainerID := newTestTrainerID(t)
	existing := newTestBooking(t, "10:00", "11:00")
	detector := domain.NewConflictDetector(
		&fakeBookingFinder{bookings: []*domain.Booking{existing}},
	)

	tests := []struct {
		name string
		from string
		to   string
		want error
	}{
		{"adjacent end", "09:00", "10:00", nil},
		{
			"partial overlap end",
			"10:30",
			"11:30",
			domain.ErrBookingConflict,
		},
		{"identical", "10:00", "11:00", domain.ErrBookingConflict},
		{
			"invalid period",
			"11:00",
			"10:00",
			domain.ErrInvalidPeriod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := detector.Detect(
				ctx,
				trainerID,
				parseHour(t, tt.from),
				parseHour(t, tt.to),
			)
			if !errors.Is(err, tt.want) {
				t.Errorf("Detect: want=%v got=%v", tt.want, err)
			}
		})
	}
}

// 別のトレーナーの予約とは衝突しない。
func TestConflictDetector_Detect_OtherTrainer(t *testing.T) {
	existing := newTestBooking(t, "10:00", "11:00")
	other, err := domain.NewID[domain.TrainerTag]("trainer2")
	if err != nil {
		t.Fatal(err)
	}
	detector := domain.NewConflictDetector(
		&fakeBookingFinder{bookings: []*domain.Booking{existing}},
	)
	if err := detector.Detect(
		context.Background(),
		other,
		parseHour(t, "10:00"),
		parseHour(t, "11:00"),
	); err != nil {
		t.Errorf("Detect for another trainer: want=nil got=%v", err)
	}
}

func TestConflictDetector_IgnoresCancelled(t *testing.T) {
	ctx := context.Background()
	trainerID := newTestTrainerID(t)
	// 既存予約をキャンセル済みとして扱う場合、衝突しないことを確認したいが
	// 本章では Booking エンティティに Cancel() を実装していない。
	// 代わりに FindByTrainerAndPeriod が空を返すケースをテストする。
	finder := &fakeBookingFinder{bookings: nil}
	detector := domain.NewConflictDetector(finder)

	if err := detector.Detect(
		ctx,
		trainerID,
		parseHour(t, "10:00"),
		parseHour(t, "11:00"),
	); err != nil {
		t.Errorf("Detect on empty: want=nil got=%v", err)
	}
}
