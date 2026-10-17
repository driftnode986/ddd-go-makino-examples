package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch06-aggregate-factory/internal/order/domain"
)

func orderID(t *testing.T, s string) domain.OrderID {
	t.Helper()
	id, err := domain.NewID[domain.OrderTag](s)
	if err != nil {
		t.Fatalf("NewID order: %v", err)
	}
	return id
}

func customerID(t *testing.T, s string) domain.CustomerID {
	t.Helper()
	id, err := domain.NewID[domain.CustomerTag](s)
	if err != nil {
		t.Fatalf("NewID customer: %v", err)
	}
	return id
}

func productID(t *testing.T, s string) domain.ProductID {
	t.Helper()
	id, err := domain.NewID[domain.ProductTag](s)
	if err != nil {
		t.Fatalf("NewID product: %v", err)
	}
	return id
}

func jpy(t *testing.T, amount int64) domain.Money {
	t.Helper()
	m, err := domain.NewMoney(amount, domain.JPY)
	if err != nil {
		t.Fatalf("NewMoney: %v", err)
	}
	return m
}

func newTestOrder(t *testing.T) *domain.Order {
	t.Helper()
	o, err := domain.NewOrder(
		orderID(t, "order001"),
		customerID(t, "cust0001"),
	)
	if err != nil {
		t.Fatalf("NewOrder: %v", err)
	}
	return o
}

// #region test-add-item-only-in-draft
func TestOrder_AddItem_OnlyInDraft(t *testing.T) {
	clock := domain.FixedClock{
		T: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	o := newTestOrder(t)

	// draft 状態で追加できる
	if err := o.AddItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	); err != nil {
		t.Fatalf("AddItem in draft: %v", err)
	}

	// 確定後は追加できない
	if err := o.Confirm(clock); err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	err := o.AddItem(
		productID(t, "prd00002"),
		1,
		jpy(t, 500),
	)
	if !errors.Is(err, domain.ErrCannotModifyConfirmed) {
		t.Errorf(
			"AddItem in confirmed: "+
				"want ErrCannotModifyConfirmed got %v",
			err,
		)
	}
}

// #endregion

func TestOrder_AddItem_RejectsInvalidQuantity(t *testing.T) {
	o := newTestOrder(t)
	err := o.AddItem(
		productID(t, "prd00001"),
		0,
		jpy(t, 1000),
	)
	if !errors.Is(err, domain.ErrInvalidQuantity) {
		t.Errorf(
			"AddItem quantity=0: want ErrInvalidQuantity got %v",
			err,
		)
	}
}

func TestOrder_RemoveItem(t *testing.T) {
	o := newTestOrder(t)
	if err := o.AddItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.AddItem(
		productID(t, "prd00002"),
		1,
		jpy(t, 500),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.RemoveItem(productID(t, "prd00001")); err != nil {
		t.Fatalf("RemoveItem: %v", err)
	}
	if len(o.Items()) != 1 {
		t.Errorf("items after remove: want 1 got %d", len(o.Items()))
	}
	err := o.RemoveItem(productID(t, "prd99999"))
	if !errors.Is(err, domain.ErrItemNotFound) {
		t.Errorf(
			"RemoveItem missing: want ErrItemNotFound got %v",
			err,
		)
	}
}

// #region test-items-returns-copy
func TestOrder_Items_ReturnsCopy(t *testing.T) {
	o := newTestOrder(t)
	if err := o.AddItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}

	items := o.Items()
	// 取得したスライスを変更しても、集約内部は変わらない
	items[0] = domain.OrderItem{} // ゼロ値で潰す

	again := o.Items()
	if again[0].Quantity() != 2 {
		t.Errorf(
			"internal state polluted: want 2 got %d",
			again[0].Quantity(),
		)
	}
}

// #endregion

// #region test-confirm-rejects-empty
func TestOrder_Confirm_RejectsEmpty(t *testing.T) {
	clock := domain.FixedClock{
		T: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	o := newTestOrder(t)
	err := o.Confirm(clock)
	if !errors.Is(err, domain.ErrEmptyOrder) {
		t.Errorf("Confirm empty: want ErrEmptyOrder got %v", err)
	}
}

// #endregion

func TestOrder_TotalAmount(t *testing.T) {
	o := newTestOrder(t)
	if err := o.AddItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := o.AddItem(
		productID(t, "prd00002"),
		3,
		jpy(t, 500),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	total, err := o.TotalAmount()
	if err != nil {
		t.Fatalf("TotalAmount: %v", err)
	}
	want := jpy(t, 2*1000+3*500)
	if !total.Equals(want) {
		t.Errorf("TotalAmount: want=%v got=%v", want, total)
	}
}
