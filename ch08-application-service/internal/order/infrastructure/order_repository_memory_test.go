package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/infrastructure"
)

func newOrder(t *testing.T, customer string) *domain.Order {
	t.Helper()
	id, err := domain.NewCustomerID(customer)
	if err != nil {
		t.Fatalf("NewCustomerID: %v", err)
	}
	return domain.NewOrder(
		id,
		time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
	)
}

func TestInMemoryOrderRepository_AddFindSave(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryOrderRepository()
	o := newOrder(t, "cust_001")

	if err := repo.Add(ctx, o); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := repo.Add(ctx, o); !errors.Is(
		err,
		domain.ErrDuplicate,
	) {
		t.Errorf("Add twice: got %v, want ErrDuplicate", err)
	}
	got, err := repo.FindByID(ctx, o.ID())
	if err != nil || got.ID() != o.ID() {
		t.Fatalf("FindByID: got %v, %v", got, err)
	}
	if err := repo.Save(ctx, o); err != nil {
		t.Errorf("Save: %v", err)
	}
	if err := repo.Save(ctx, newOrder(t, "cust_001")); !errors.Is(
		err,
		domain.ErrNotFound,
	) {
		t.Errorf("Save unknown: got %v, want ErrNotFound", err)
	}
	missing, err := domain.NewOrderID("ord_missing")
	if err != nil {
		t.Fatalf("NewOrderID: %v", err)
	}
	if _, err := repo.FindByID(ctx, missing); !errors.Is(
		err,
		domain.ErrNotFound,
	) {
		t.Errorf("FindByID unknown: got %v, want ErrNotFound", err)
	}
}

func TestInMemoryOrderRepository_ListByCustomer(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewInMemoryOrderRepository()
	for _, c := range []string{"cust_001", "cust_001", "cust_002"} {
		if err := repo.Add(ctx, newOrder(t, c)); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}
	id, err := domain.NewCustomerID("cust_001")
	if err != nil {
		t.Fatalf("NewCustomerID: %v", err)
	}
	got, err := repo.ListByCustomer(ctx, id)
	if err != nil {
		t.Fatalf("ListByCustomer: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d orders, want 2", len(got))
	}
}

func TestStubPricing_UnitPrice(t *testing.T) {
	p := infrastructure.NewStubPricing()
	id, err := domain.NewProductID("prod_001")
	if err != nil {
		t.Fatalf("NewProductID: %v", err)
	}
	m, err := p.UnitPrice(context.Background(), id)
	if err != nil {
		t.Fatalf("UnitPrice: %v", err)
	}
	if m.Amount() != 1000 || m.Currency() != "JPY" {
		t.Errorf(
			"got %d %s, want 1000 JPY",
			m.Amount(),
			m.Currency(),
		)
	}
}
