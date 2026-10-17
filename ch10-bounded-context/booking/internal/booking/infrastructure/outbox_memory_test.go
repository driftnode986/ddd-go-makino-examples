package infrastructure_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/infrastructure"
)

func confirmedOrder(t *testing.T) *domain.Order {
	t.Helper()
	customerID, err := domain.NewCustomerID("cus_001")
	if err != nil {
		t.Fatal(err)
	}
	o := domain.NewOrder(customerID, time.Unix(1700000000, 0))
	pid, _ := domain.NewProductID("prod_001")
	price, _ := domain.NewMoney(1500, "JPY")
	if err := o.AddItem(pid, 2, price); err != nil {
		t.Fatal(err)
	}
	if err := o.Confirm(time.Unix(1700000100, 0)); err != nil {
		t.Fatal(err)
	}
	return o
}

func TestMemoryOutbox_AppendAndPublish(t *testing.T) {
	o := confirmedOrder(t)
	outbox := infrastructure.NewMemoryOutbox()
	if err := outbox.Append(context.Background(), o.PullEvents()); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries := outbox.Unpublished()
	if len(entries) != 1 {
		t.Fatalf("unpublished = %d, want 1", len(entries))
	}
	got := entries[0]
	if got.ID != 1 || got.EventName != "booking.OrderPlaced" || got.AggregateID != o.ID().String() {
		t.Errorf("entry = %+v", got)
	}
	var payload map[string]any
	if err := json.Unmarshal(got.Payload, &payload); err != nil {
		t.Fatalf("payload is not JSON: %v", err)
	}
	if payload["order_id"] != o.ID().String() || payload["total_amount"] != float64(3000) {
		t.Errorf("payload = %v", payload)
	}

	outbox.MarkPublished(got.ID, time.Now())
	if n := len(outbox.Unpublished()); n != 0 {
		t.Errorf("unpublished after MarkPublished = %d, want 0", n)
	}
}

type unknownEvent struct{}

func (unknownEvent) EventName() string     { return "booking.Unknown" }
func (unknownEvent) OccurredAt() time.Time { return time.Time{} }
func (unknownEvent) AggregateID() string   { return "x" }

func TestMemoryOutbox_RejectsUnknownEvent(t *testing.T) {
	outbox := infrastructure.NewMemoryOutbox()
	err := outbox.Append(context.Background(), []domain.DomainEvent{unknownEvent{}})
	if err == nil {
		t.Fatal("Append(unknown event): want error, got nil")
	}
	if n := len(outbox.Unpublished()); n != 0 {
		t.Errorf("unpublished = %d, want 0", n)
	}
}

func TestMemoryOrderRepository_AddAndFind(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryOrderRepository()
	o := confirmedOrder(t)
	if err := repo.Add(ctx, o); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindByID(ctx, o.ID())
	if err != nil {
		t.Fatal(err)
	}
	if got.Status() != domain.StatusConfirmed {
		t.Errorf("status = %v, want confirmed", got.Status())
	}
	list, err := repo.ListByCustomer(ctx, o.CustomerID())
	if err != nil || len(list) != 1 {
		t.Errorf("ListByCustomer = %d, %v", len(list), err)
	}
}
