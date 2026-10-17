package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
)

func mustCustomerID(t *testing.T, s string) domain.CustomerID {
	t.Helper()
	id, err := domain.NewCustomerID(s)
	if err != nil {
		t.Fatalf("NewCustomerID(%q): %v", s, err)
	}
	return id
}

func TestNewOrder_StartsAsDraft(t *testing.T) {
	now := time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC)
	o, err := domain.NewOrder(mustCustomerID(t, "cus_001"), now)
	if err != nil {
		t.Fatalf("NewOrder: %v", err)
	}
	if o.Status() != domain.StatusDraft {
		t.Errorf("status = %q, want %q", o.Status(), domain.StatusDraft)
	}
	if !o.CreatedAt().Equal(now) {
		t.Errorf("createdAt = %v, want %v", o.CreatedAt(), now)
	}
	if o.ID().String() == "" {
		t.Error("order id is empty")
	}
}

func TestNewOrder_RequiresCustomerID(t *testing.T) {
	if _, err := domain.NewOrder(domain.CustomerID{}, time.Now()); err == nil {
		t.Fatal("NewOrder with zero customer id: want error, got nil")
	}
}

func TestOrder_AddItemAndTotal(t *testing.T) {
	o, err := domain.NewOrder(mustCustomerID(t, "cus_001"), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	pid, _ := domain.NewProductID("prod_001")
	if err := o.AddItem(pid, 2, domain.MustMoney(1500, "JPY")); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.AddItem(pid, 1, domain.MustMoney(500, "JPY")); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if got := o.Total().Amount(); got != 3500 {
		t.Errorf("total = %d, want 3500", got)
	}
	if got := len(o.Items()); got != 2 {
		t.Errorf("items = %d, want 2", got)
	}
}

func TestOrder_AddItem_InvalidQuantity(t *testing.T) {
	o, _ := domain.NewOrder(mustCustomerID(t, "cus_001"), time.Now())
	pid, _ := domain.NewProductID("prod_001")
	err := o.AddItem(pid, 0, domain.MustMoney(1500, "JPY"))
	if !errors.Is(err, domain.ErrInvalidQuantity) {
		t.Fatalf("err = %v, want ErrInvalidQuantity", err)
	}
}

func TestOrder_ConfirmedOrderCannotBeModified(t *testing.T) {
	o, _ := domain.NewOrder(mustCustomerID(t, "cus_001"), time.Now())
	if err := o.Confirm(); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if o.Status() != domain.StatusConfirmed {
		t.Errorf("status = %q, want %q", o.Status(), domain.StatusConfirmed)
	}
	pid, _ := domain.NewProductID("prod_001")
	err := o.AddItem(pid, 1, domain.MustMoney(100, "JPY"))
	if !errors.Is(err, domain.ErrCannotModify) {
		t.Fatalf("err = %v, want ErrCannotModify", err)
	}
	if err := o.Confirm(); !errors.Is(err, domain.ErrCannotModify) {
		t.Fatalf("second Confirm err = %v, want ErrCannotModify", err)
	}
}
