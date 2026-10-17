package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

func mustCustomer(t *testing.T) domain.CustomerID {
	t.Helper()
	id, err := domain.NewCustomerID("cust_001")
	if err != nil {
		t.Fatalf("NewCustomerID: %v", err)
	}
	return id
}

func mustProduct(t *testing.T, s string) domain.ProductID {
	t.Helper()
	id, err := domain.NewProductID(s)
	if err != nil {
		t.Fatalf("NewProductID: %v", err)
	}
	return id
}

func TestNewOrder_Draft(t *testing.T) {
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	o := domain.NewOrder(mustCustomer(t), at)
	if o.Status() != domain.StatusDraft {
		t.Errorf("Status = %v, want draft", o.Status())
	}
	if o.ID().String() == "" {
		t.Error("ID is empty")
	}
	if !o.PlacedAt().Equal(at) {
		t.Errorf("PlacedAt = %v, want %v", o.PlacedAt(), at)
	}
}

func TestOrder_AddItemAndTotal(t *testing.T) {
	o := domain.NewOrder(mustCustomer(t), time.Time{})
	if got := o.Total(); got.Amount() != 0 {
		t.Errorf("Total of empty order = %d, want 0", got.Amount())
	}
	if err := o.AddItem(
		mustProduct(t, "prod_001"),
		2,
		domain.MustMoney(1000, "JPY"),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.AddItem(
		mustProduct(t, "prod_002"),
		3,
		domain.MustMoney(500, "JPY"),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if got := o.Total().Amount(); got != 3500 {
		t.Errorf("Total = %d, want 3500", got)
	}
	if n := len(o.Items()); n != 2 {
		t.Errorf("Items = %d, want 2", n)
	}
}

func TestOrder_AddItem_InvalidQuantity(t *testing.T) {
	o := domain.NewOrder(mustCustomer(t), time.Time{})
	err := o.AddItem(
		mustProduct(t, "prod_001"),
		0,
		domain.MustMoney(1000, "JPY"),
	)
	if !errors.Is(err, domain.ErrInvalidQuantity) {
		t.Errorf("got %v, want ErrInvalidQuantity", err)
	}
}

func TestOrder_Confirm_RejectsEmpty(t *testing.T) {
	o := domain.NewOrder(mustCustomer(t), time.Time{})
	if err := o.Confirm(); !errors.Is(err, domain.ErrEmptyOrder) {
		t.Errorf("Confirm empty: got %v, want ErrEmptyOrder", err)
	}
	if o.Status() != domain.StatusDraft {
		t.Errorf("Status = %v, want draft", o.Status())
	}
}

func TestOrder_Confirm_ThenCannotModify(t *testing.T) {
	o := domain.NewOrder(mustCustomer(t), time.Time{})
	if err := o.AddItem(
		mustProduct(t, "prod_001"),
		1,
		domain.MustMoney(1000, "JPY"),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.Confirm(); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if err := o.Confirm(); !errors.Is(err, domain.ErrCannotModify) {
		t.Errorf("Confirm twice: got %v, want ErrCannotModify", err)
	}
	err := o.AddItem(
		mustProduct(t, "prod_001"),
		1,
		domain.MustMoney(1000, "JPY"),
	)
	if !errors.Is(err, domain.ErrCannotModify) {
		t.Errorf(
			"AddItem after Confirm: got %v, want ErrCannotModify",
			err,
		)
	}
}

func TestNewIDs_RejectEmpty(t *testing.T) {
	if _, err := domain.NewCustomerID(""); !errors.Is(
		err,
		domain.ErrInvalidID,
	) {
		t.Errorf("NewCustomerID: got %v", err)
	}
	if _, err := domain.NewProductID(""); !errors.Is(
		err,
		domain.ErrInvalidID,
	) {
		t.Errorf("NewProductID: got %v", err)
	}
	if _, err := domain.NewOrderID(""); !errors.Is(
		err,
		domain.ErrInvalidID,
	) {
		t.Errorf("NewOrderID: got %v", err)
	}
}

func TestMoney_Add_CurrencyMismatch(t *testing.T) {
	_, err := domain.MustMoney(100, "JPY").
		Add(domain.MustMoney(1, "USD"))
	if !errors.Is(err, domain.ErrCurrencyMismatch) {
		t.Errorf("got %v, want ErrCurrencyMismatch", err)
	}
}
