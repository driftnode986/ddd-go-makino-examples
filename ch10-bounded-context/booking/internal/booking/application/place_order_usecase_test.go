package application_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/infrastructure"
)

func TestPlaceOrderUseCase_EmitsOrderPlacedToOutbox(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	outbox := infrastructure.NewMemoryOutbox()
	uc := &application.PlaceOrderUseCase{
		OrderRepo: repo,
		Outbox:    outbox,
		Clock:     application.FixedClock{At: time.Unix(1700000000, 0)},
		Logger:    slog.Default(),
	}

	cmd := application.PlaceOrderCommand{
		CustomerID: "cus_001",
		Items: []application.OrderItemInput{
			{ProductID: "prod_001", Quantity: 2, UnitPrice: 1500},
		},
	}

	result, err := uc.Execute(context.Background(), cmd)
	if err != nil {
		t.Fatal(err)
	}
	if result.OrderID == "" {
		t.Fatal("order id is empty")
	}
	if result.Total != 3000 {
		t.Errorf("total = %d, want 3000", result.Total)
	}

	entries := outbox.Unpublished()
	if len(entries) != 1 {
		t.Fatalf("outbox entries = %d, want 1", len(entries))
	}
	if entries[0].EventName != "booking.OrderPlaced" {
		t.Errorf("event name = %s, want booking.OrderPlaced", entries[0].EventName)
	}
}
