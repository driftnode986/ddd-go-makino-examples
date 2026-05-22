package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure"
)

func mustBookingID(t *testing.T, s string) domain.BookingID {
	t.Helper()
	id, err := domain.NewBookingID(s)
	if err != nil {
		t.Fatalf("NewBookingID(%q): %v", s, err)
	}
	return id
}

func mustMemberID(t *testing.T, s string) domain.MemberID {
	t.Helper()
	id, err := domain.NewMemberID(s)
	if err != nil {
		t.Fatalf("NewMemberID(%q): %v", s, err)
	}
	return id
}

func mustLessonID(t *testing.T, s string) domain.LessonID {
	t.Helper()
	id, err := domain.NewLessonID(s)
	if err != nil {
		t.Fatalf("NewLessonID(%q): %v", s, err)
	}
	return id
}

func mustNewBooking(t *testing.T, isoTime string) *domain.Booking {
	t.Helper()
	at, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		t.Fatalf("Parse(%q): %v", isoTime, err)
	}
	return domain.NewBooking(
		mustBookingID(t, "bk_001"),
		mustMemberID(t, "mem_001"),
		mustLessonID(t, "lsn_001"),
		at,
	)
}

func mustNewBookingWithVersion(t *testing.T, id domain.BookingID, version int) *domain.Booking {
	t.Helper()
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	return domain.ReconstructBooking(
		id,
		mustMemberID(t, "mem_001"),
		mustLessonID(t, "lsn_001"),
		domain.StatusConfirmed,
		at,
		version,
	)
}

func TestInMemoryBookingRepository_AddAndFind(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")

	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}
	got, err := repo.FindByID(ctx, b.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if diff := cmp.Diff(
		b, got, cmp.AllowUnexported(domain.Booking{}, domain.BookingID{}, domain.MemberID{}, domain.LessonID{}),
	); diff != "" {
		t.Errorf("Booking mismatch (-want +got):\n%s", diff)
	}
}

func TestInMemoryBookingRepository_AddDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")

	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add #1: %v", err)
	}
	err := repo.Add(ctx, b)
	if !errors.Is(err, domain.ErrDuplicate) {
		t.Errorf("got %v, want ErrDuplicate", err)
	}
}

func TestInMemoryBookingRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()

	_, err := repo.FindByID(ctx, mustBookingID(t, "missing"))
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}

func TestInMemoryBookingRepository_Save_VersionConflict(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Add 後の保存済み version は 1。version=5 で Save しようとすると競合になる。
	stale := mustNewBookingWithVersion(t, b.ID(), 5)
	err := repo.Save(ctx, stale)
	if !errors.Is(err, domain.ErrVersionConflict) {
		t.Errorf("got %v, want ErrVersionConflict", err)
	}
}

func TestInMemoryBookingRepository_Save_HappyPath(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// 同じ Booking を Confirm して Save (version 1 → 2)
	if err := b.Confirm(); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if err := repo.Save(ctx, b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.FindByID(ctx, b.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.Status() != domain.StatusConfirmed {
		t.Errorf("status = %v, want %v", got.Status(), domain.StatusConfirmed)
	}
	if got.Version() != 2 {
		t.Errorf("version = %d, want 2", got.Version())
	}
}

func TestInMemoryBookingRepository_Remove(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if err := repo.Remove(ctx, b.ID()); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := repo.FindByID(ctx, b.ID())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("after Remove: got %v, want ErrNotFound", err)
	}
}

func TestInMemoryBookingRepository_Remove_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()

	err := repo.Remove(ctx, mustBookingID(t, "missing"))
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}
