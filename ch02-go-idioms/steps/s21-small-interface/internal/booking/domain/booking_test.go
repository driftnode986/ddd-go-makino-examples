package domain

import (
	"context"
	"testing"
	"time"
)

// memoryRepo は 3 メソッドの小さい interface を満たす最小の実装。
type memoryRepo struct{ m map[BookingID]*Booking }

func (r *memoryRepo) Save(
	_ context.Context,
	b *Booking,
) error {
	r.m[b.ID()] = b
	return nil
}

func (r *memoryRepo) FindByID(
	_ context.Context,
	id BookingID,
) (*Booking, error) {
	return r.m[id], nil
}

func (r *memoryRepo) Remove(
	_ context.Context,
	id BookingID,
) error {
	delete(r.m, id)
	return nil
}

var _ BookingRepository = (*memoryRepo)(nil)

func TestBookingRepository_SaveFindRemove(t *testing.T) {
	ctx := context.Background()
	repo := &memoryRepo{m: map[BookingID]*Booking{}}
	b := &Booking{status: "pending"}
	if err := repo.Save(ctx, b); err != nil {
		t.Fatal(err)
	}
	if got, _ := repo.FindByID(ctx, b.ID()); got != b {
		t.Fatal("FindByID should return the saved booking")
	}
	if err := repo.Remove(ctx, b.ID()); err != nil {
		t.Fatal(err)
	}
	if got, _ := repo.FindByID(ctx, b.ID()); got != nil {
		t.Fatal("Remove should delete the booking")
	}
}

// 中身を省略した版なので、状態変更メソッドは何もせず nil を返す。
func TestBooking_StubMethodsKeepState(t *testing.T) {
	b := &Booking{status: "pending"}
	for name, call := range map[string]func() error{
		"Confirm": b.Confirm,
		"Cancel":  func() error { return b.Cancel("x", time.Now) },
		"Reschedule": func() error {
			return b.Reschedule(time.Now())
		},
		"MarkCompleted": b.MarkCompleted,
	} {
		if err := call(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	if b.Status() != "pending" {
		t.Fatalf("Status: got %s", b.Status())
	}
}
