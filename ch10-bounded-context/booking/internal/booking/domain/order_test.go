// #region order-test-imports
package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
)

// #endregion

// #region confirm-emits-order-placed
func TestOrder_Confirm_EmitsOrderPlaced(t *testing.T) {
	customerID, _ := domain.NewCustomerID("cus_001")
	o := domain.NewOrder(customerID, time.Unix(1700000000, 0))
	pid, _ := domain.NewProductID("prod_001")
	price, _ := domain.NewMoney(1500, "JPY")
	_ = o.AddItem(pid, 2, price)

	if err := o.Confirm(time.Unix(1700000100, 0)); err != nil {
		t.Fatalf("confirm: %v", err)
	}

	events := o.PullEvents()
	if len(events) != 1 {
		t.Fatalf("event count = %d, want 1", len(events))
	}
	if events[0].EventName() != "booking.OrderPlaced" {
		t.Errorf("event name = %s, want booking.OrderPlaced",
			events[0].EventName())
	}

	// 2 回目の PullEvents は空であること
	if got := o.PullEvents(); len(got) != 0 {
		t.Errorf("second pull should be empty, got %d", len(got))
	}
}

// #endregion

func TestOrder_Confirm_Twice_Fails(t *testing.T) {
	customerID, _ := domain.NewCustomerID("cus_001")
	o := domain.NewOrder(customerID, time.Unix(1700000000, 0))
	pid, _ := domain.NewProductID("prod_001")
	price, _ := domain.NewMoney(1500, "JPY")
	_ = o.AddItem(pid, 1, price)

	if err := o.Confirm(time.Unix(1700000100, 0)); err != nil {
		t.Fatal(err)
	}
	if err := o.Confirm(time.Unix(1700000200, 0)); err == nil {
		t.Fatal("expected error on second confirm, got nil")
	}
}

func TestOrderPlaced_SnapshotsOrder(t *testing.T) {
	customerID, _ := domain.NewCustomerID("cus_001")
	o := domain.NewOrder(customerID, time.Unix(1700000000, 0))
	pid, _ := domain.NewProductID("prod_001")
	price, _ := domain.NewMoney(1500, "JPY")
	if err := o.AddItem(pid, 2, price); err != nil {
		t.Fatal(err)
	}
	at := time.Unix(1700000100, 0)
	e := domain.NewOrderPlaced(o, at)

	if e.OrderID() != o.ID() || e.AggregateID() != o.ID().String() {
		t.Errorf(
			"order id = %v / %s, want %v",
			e.OrderID(),
			e.AggregateID(),
			o.ID(),
		)
	}
	if e.CustomerID() != customerID || e.TotalAmount() != 3000 {
		t.Errorf(
			"customer = %v, total = %d",
			e.CustomerID(),
			e.TotalAmount(),
		)
	}
	if !e.OccurredAt().Equal(at) {
		t.Errorf("occurredAt = %v, want %v", e.OccurredAt(), at)
	}
}
