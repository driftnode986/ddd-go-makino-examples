package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

func TestPlaceOrder(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	items := []domain.OrderItem{
		{ProductID: "p1", Quantity: 2, UnitPrice: 1500},
	}
	o, err := domain.PlaceOrder(items, "cus_1", now)
	if err != nil {
		t.Fatal(err)
	}
	if o.Status() != domain.OrderStatusPlaced ||
		o.Subtotal() != 3000 ||
		!o.PlacedAt().Equal(now) {
		t.Errorf(
			"order = %v %d %v",
			o.Status(),
			o.Subtotal(),
			o.PlacedAt(),
		)
	}
	items[0].Quantity = 99 // 呼び出し側のスライスを書き換えても注文は変わらない
	if o.Items()[0].Quantity != 2 {
		t.Error("PlaceOrder must copy the items")
	}

	for name, in := range map[string][]domain.OrderItem{
		"no items": nil,
		"zero quantity": {{
			ProductID: "p1",
			Quantity:  0,
			UnitPrice: 100,
		}},
	} {
		if _, err := domain.PlaceOrder(in, "cus_1", now); !errors.Is(
			err,
			domain.ErrInvalidOrder,
		) {
			t.Errorf("%s: err = %v, want ErrInvalidOrder", name, err)
		}
	}
}

func TestOrder_Cancel(t *testing.T) {
	o, _ := domain.PlaceOrder(
		[]domain.OrderItem{
			{ProductID: "p1", Quantity: 1, UnitPrice: 100},
		},
		"cus_1",
		time.Now(),
	)
	if err := o.Cancel(); err != nil {
		t.Fatal(err)
	}
	if err := o.Cancel(); !errors.Is(err, domain.ErrInvalidStatus) {
		t.Errorf(
			"second Cancel err = %v, want ErrInvalidStatus",
			err,
		)
	}
}

func TestPricingService_Total(t *testing.T) {
	o, _ := domain.PlaceOrder(
		[]domain.OrderItem{
			{ProductID: "p1", Quantity: 3, UnitPrice: 999},
		},
		"cus_1",
		time.Now(),
	)
	// 2997 + 299（1 円未満切り捨て）
	if got := domain.NewPricingService(10).Total(o); got != 3296 {
		t.Errorf("Total = %d, want 3296", got)
	}
}

func TestOrderPlaced_EventName(t *testing.T) {
	var e domain.DomainEvent = domain.OrderPlaced{OrderID: "ord_1"}
	if e.EventName() != "booking.OrderPlaced" {
		t.Errorf("EventName = %q", e.EventName())
	}
}
