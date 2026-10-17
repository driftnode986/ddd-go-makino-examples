package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/infrastructure/memory"
)

func TestBookingRepository_SaveFindList(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewBookingRepository()
	member := domain.NewMemberID()
	b := domain.NewBooking(member, domain.NewLessonID(), time.Now())

	if err := repo.Save(ctx, b); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := repo.FindByID(ctx, b.ID())
	if err != nil || got.ID() != b.ID() {
		t.Fatalf("FindByID = %v, %v", got, err)
	}
	list, err := repo.ListByMember(ctx, member)
	if err != nil || len(list) != 1 {
		t.Errorf("ListByMember = %d, %v", len(list), err)
	}
	if list, _ := repo.ListByMember(ctx, domain.NewMemberID()); len(list) != 0 {
		t.Errorf("other member's bookings = %d, want 0", len(list))
	}
}

func TestBookingRepository_NotFound(t *testing.T) {
	repo := memory.NewBookingRepository()
	_, err := repo.FindByID(context.Background(), domain.NewBookingID())
	if !errors.Is(err, domain.ErrBookingNotFound) {
		t.Fatalf("err = %v, want ErrBookingNotFound", err)
	}
}
