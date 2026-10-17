package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch04-entity/internal/booking/domain"
)

// helper: テスト用の Booking を組み立てる。
func newTestBooking(t *testing.T, clock func() time.Time) *domain.Booking {
	t.Helper()
	bid, _ := domain.NewID[domain.BookingTag]("book-0001")
	mid, _ := domain.NewID[domain.MemberTag]("memb-0001")
	lid, _ := domain.NewID[domain.LessonTag]("less-0001")
	scheduled := clock().Add(24 * time.Hour)
	b, err := domain.NewBooking(bid, mid, lid, scheduled, clock)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	return b
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNewBooking_RejectsScheduledInPast(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)

	bid, _ := domain.NewID[domain.BookingTag]("book-0001")
	mid, _ := domain.NewID[domain.MemberTag]("memb-0001")
	lid, _ := domain.NewID[domain.LessonTag]("less-0001")
	past := now.Add(-1 * time.Hour)

	_, err := domain.NewBooking(bid, mid, lid, past, clock)
	if !errors.Is(err, domain.ErrScheduledInPast) {
		t.Fatalf("want ErrScheduledInPast, got %v", err)
	}
}

func TestBooking_Confirm_PendingToConfirmed(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	b := newTestBooking(t, clock)

	if got := b.Status(); got != domain.StatusPending {
		t.Fatalf("initial status: want pending, got %v", got)
	}
	if err := b.Confirm(clock); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if got := b.Status(); got != domain.StatusConfirmed {
		t.Fatalf("after Confirm: want confirmed, got %v", got)
	}
	if b.Version() != 2 {
		t.Fatalf("version after Confirm: want 2, got %d", b.Version())
	}
}

func TestBooking_Confirm_RejectsFromCancelled(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	b := newTestBooking(t, clock)

	if err := b.Cancel(clock); err != nil {
		t.Fatal(err)
	}
	if err := b.Confirm(clock); !errors.Is(err, domain.ErrCannotConfirm) {
		t.Fatalf("Confirm from cancelled: want ErrCannotConfirm, got %v", err)
	}
}

// #region cancel-table
func TestBooking_Cancel_FromPendingOrConfirmed(t *testing.T) {
	tests := []struct {
		name    string
		actions []func(*domain.Booking, func() time.Time) error
		want    error
	}{
		{
			"cancel from pending",
			nil,
			nil,
		},
		{
			"cancel from confirmed",
			[]func(*domain.Booking, func() time.Time) error{
				(*domain.Booking).Confirm,
			},
			nil,
		},
		{
			"cannot cancel from completed",
			[]func(*domain.Booking, func() time.Time) error{
				(*domain.Booking).Confirm,
				(*domain.Booking).Complete,
			},
			domain.ErrCannotCancel,
		},
	}
	// #endregion
	// #region cancel-run
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTestBooking(t, clock)
			for _, act := range tt.actions {
				if err := act(b, clock); err != nil {
					t.Fatalf("setup action: %v", err)
				}
			}
			got := b.Cancel(clock)
			if !errors.Is(got, tt.want) {
				t.Fatalf("Cancel: want=%v got=%v", tt.want, got)
			}
		})
	}
}

// #endregion

func TestBooking_Complete_OnlyFromConfirmed(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	b := newTestBooking(t, clock)

	if err := b.Complete(clock); !errors.Is(err, domain.ErrCannotComplete) {
		t.Fatalf("Complete from pending: want ErrCannotComplete, got %v", err)
	}
	if err := b.Confirm(clock); err != nil {
		t.Fatal(err)
	}
	if err := b.Complete(clock); err != nil {
		t.Fatalf("Complete from confirmed: %v", err)
	}
	if b.Status() != domain.StatusCompleted {
		t.Fatalf("after Complete: got %v", b.Status())
	}
}

// #region identity
func TestBooking_IdentityIsByID(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	a := newTestBooking(t, clock)
	b := newTestBooking(t, clock) // 同じ ID で再構築

	// エンティティの等価性は ID で判定する (値そのものではない)
	if !a.ID().Equals(b.ID()) {
		t.Errorf("same ID should be equal")
	}
	// 状態を変えた a と、変えていない b は内容が違う
	if err := a.Confirm(clock); err != nil {
		t.Fatal(err)
	}
	if a.Status() == b.Status() {
		t.Errorf("a was confirmed, b should remain pending")
	}
	// それでも ID は等しい (エンティティとしては同一)
	if !a.ID().Equals(b.ID()) {
		t.Errorf("identity must be preserved across state changes")
	}
}

// #endregion
