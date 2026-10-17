package infrastructure_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure"
)

// 本に載せた2件 (booking_repository_memory_test.go) のほかの振る舞いを確かめる。

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

// 取り出した集約を変更して Save する、ユースケースと同じ流れ。
// 保存済みの集約とポインタを共有していると、ここで競合と誤判定される。
func TestInMemoryBookingRepository_Save_AfterFindByID(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	if err := repo.Add(ctx, mustNewBooking(t, "2026-06-01T10:00:00Z")); err != nil {
		t.Fatalf("Add: %v", err)
	}
	b, err := repo.FindByID(ctx, mustBookingID(t, "bk_001"))
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if err := b.Confirm(); err != nil { // version 1 → 2
		t.Fatalf("Confirm: %v", err)
	}
	if err := repo.Save(ctx, b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.FindByID(ctx, b.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.Status() != domain.StatusConfirmed || got.Version() != 2 {
		t.Errorf("got status=%v version=%d, want confirmed/2", got.Status(), got.Version())
	}
}

// Save した後に呼び出し側の集約を変更しても、保存済みの状態は変わらない。
func TestInMemoryBookingRepository_StoresCopy(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()
	b := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, b); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := b.Cancel(); err != nil { // 保存せずに変更だけする
		t.Fatalf("Cancel: %v", err)
	}

	got, err := repo.FindByID(ctx, b.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.Status() != domain.StatusPending || got.Version() != 1 {
		t.Errorf("got status=%v version=%d, want pending/1", got.Status(), got.Version())
	}
}

func TestInMemoryBookingRepository_Save_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryBookingRepository()

	err := repo.Save(ctx, mustNewBooking(t, "2026-06-01T10:00:00Z"))
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
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
	if err := repo.Remove(ctx, b.ID()); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Remove twice: got %v, want ErrNotFound", err)
	}
}
