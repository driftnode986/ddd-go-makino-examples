package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/infrastructure"
)

func newOrder(t *testing.T, customer string) *domain.Order {
	t.Helper()
	cid, err := domain.NewCustomerID(customer)
	if err != nil {
		t.Fatal(err)
	}
	o, err := domain.NewOrder(cid, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func TestMemoryOrderRepository_AddAndFind(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryOrderRepository()
	o := newOrder(t, "cus_001")

	if err := repo.Add(ctx, o); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := repo.Add(ctx, o); !errors.Is(err, domain.ErrDuplicate) {
		t.Fatalf("second Add err = %v, want ErrDuplicate", err)
	}
	got, err := repo.FindByID(ctx, o.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got.ID() != o.ID() {
		t.Errorf("id = %v, want %v", got.ID(), o.ID())
	}
}

func TestMemoryOrderRepository_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryOrderRepository()
	o := newOrder(t, "cus_001")

	if _, err := repo.FindByID(ctx, o.ID()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("FindByID err = %v, want ErrNotFound", err)
	}
	if err := repo.Save(ctx, o); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Save err = %v, want ErrNotFound", err)
	}
}

func TestMemoryOrderRepository_ListByCustomer(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryOrderRepository()
	a1, a2, b := newOrder(t, "cus_a"), newOrder(t, "cus_a"), newOrder(t, "cus_b")
	for _, o := range []*domain.Order{a1, a2, b} {
		if err := repo.Add(ctx, o); err != nil {
			t.Fatal(err)
		}
	}
	got, err := repo.ListByCustomer(ctx, a1.CustomerID())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("orders for cus_a = %d, want 2", len(got))
	}
}
