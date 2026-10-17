// #region memory-test-imports
package infrastructure_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure"
)

// #endregion

// #region memory-test-add-find
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
		b, got, cmp.AllowUnexported(
			domain.Booking{}, domain.BookingID{},
			domain.MemberID{}, domain.LessonID{},
		),
	); diff != "" {
		t.Errorf("Booking mismatch (-want +got):\n%s", diff)
	}
}

// #endregion

// #region memory-test-version-conflict
func TestInMemoryBookingRepository_VersionConflict(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// バージョン不整合な集約を Save しようとして競合エラーになる
	stale := mustNewBookingWithVersion(t, b.ID(), 5)
	err := repo.Save(ctx, stale)
	if err == nil {
		t.Fatal("expected ErrVersionConflict, got nil")
	}
	if want := domain.ErrVersionConflict; err != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

// #endregion
