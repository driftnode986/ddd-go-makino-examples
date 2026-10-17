package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch06-aggregate-factory/internal/order/domain"
)

type stubIDGen struct {
	id domain.OrderID
}

func (s *stubIDGen) NewOrderID() domain.OrderID { return s.id }

func TestNewOrder_RejectsZero(t *testing.T) {
	_, err := domain.NewOrder(
		domain.OrderID{},
		customerID(t, "cust0001"),
	)
	if err == nil {
		t.Errorf("NewOrder with zero id: want error got nil")
	}
}

func TestOrderFactory_Create(t *testing.T) {
	id := orderID(t, "order001")
	factory := domain.NewOrderFactory(&stubIDGen{id: id})
	customer := customerID(t, "cust0001")
	o, err := factory.Create(customer)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !o.CustomerID().Equals(customer) {
		t.Errorf(
			"CustomerID: want=%v got=%v", customer, o.CustomerID(),
		)
	}
	if !o.ID().Equals(id) {
		t.Errorf("ID: want=%v got=%v", id, o.ID())
	}
	if o.Status() != domain.StatusDraft {
		t.Errorf("Status: want=draft got=%v", o.Status())
	}
}

// #region test-with-options
func TestNewOrderWithOptions_AppliesOptions(t *testing.T) {
	id := orderID(t, "order001")
	customer := customerID(t, "cust0001")
	rate, err := domain.NewDiscountRate(0.1)
	if err != nil {
		t.Fatalf("NewDiscountRate: %v", err)
	}

	o, err := domain.NewOrderWithOptions(id, customer,
		domain.WithNotes("ギフト包装"),
		domain.WithDiscountRate(rate),
	)
	if err != nil {
		t.Fatalf("NewOrderWithOptions: %v", err)
	}
	if o.Notes() != "ギフト包装" {
		t.Errorf("Notes: want=ギフト包装 got=%s", o.Notes())
	}
	if o.DiscountRate().Value() != 0.1 {
		t.Errorf(
			"DiscountRate: want=0.1 got=%v",
			o.DiscountRate().Value(),
		)
	}
}

// #endregion

func TestNewOrderWithOptions_ValidatesAfterOptions(t *testing.T) {
	_, err := domain.NewOrderWithOptions(
		domain.OrderID{},
		customerID(t, "cust0001"),
		domain.WithNotes("test"),
	)
	if err == nil {
		t.Errorf(
			"NewOrderWithOptions with zero id: want error got nil",
		)
	}
}

func TestNewOrder_AndReconstructOrder_PreserveItems(t *testing.T) {
	o, err := domain.NewOrder(
		orderID(t, "order001"),
		customerID(t, "cust0001"),
	)
	if err != nil {
		t.Fatalf("NewOrder: %v", err)
	}
	if err := o.AddItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	items := o.Items()
	placedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	restored := domain.ReconstructOrder(
		o.ID(),
		o.CustomerID(),
		items,
		domain.StatusConfirmed,
		placedAt,
	)
	if !restored.ID().Equals(o.ID()) {
		t.Errorf("ID mismatch after reconstruct")
	}
	if restored.Status() != domain.StatusConfirmed {
		t.Errorf("Status: want=confirmed got=%v", restored.Status())
	}
	if !restored.PlacedAt().Equal(placedAt) {
		t.Errorf(
			"PlacedAt: want=%v got=%v",
			placedAt,
			restored.PlacedAt(),
		)
	}
	if len(restored.Items()) != 1 {
		t.Errorf("Items: want 1 got %d", len(restored.Items()))
	}
}

func TestReconstructOrder_DefensiveCopy(t *testing.T) {
	id := orderID(t, "order001")
	customer := customerID(t, "cust0001")
	item, err := domain.NewOrderItem(
		productID(t, "prd00001"),
		2,
		jpy(t, 1000),
	)
	if err != nil {
		t.Fatalf("NewOrderItem: %v", err)
	}
	items := []domain.OrderItem{item}
	placedAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	restored := domain.ReconstructOrder(
		id,
		customer,
		items,
		domain.StatusConfirmed,
		placedAt,
	)
	// 渡したスライスを外部から潰しても集約には影響しないこと
	items[0] = domain.OrderItem{}
	if restored.Items()[0].Quantity() != 2 {
		t.Errorf("reconstructed order polluted by external mutation")
	}
}

func TestNewDiscountRate_Invalid(t *testing.T) {
	for _, v := range []float64{-0.1, 1.1, 2.0} {
		if _, err := domain.NewDiscountRate(v); !errors.Is(
			err,
			domain.ErrInvalidDiscountRate,
		) {
			t.Errorf(
				"NewDiscountRate(%v): "+
					"want ErrInvalidDiscountRate got %v",
				v,
				err,
			)
		}
	}
}
