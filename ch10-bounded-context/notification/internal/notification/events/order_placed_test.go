package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/events"
)

// Booking BC の outbox が書く JSON（infrastructure/outbox_memory.go の marshalEvent）を
// Notification BC の型で読めることを確かめる。
func TestOrderPlaced_DecodesPublishedPayload(t *testing.T) {
	payload := `{"order_id":"ord_001","customer_id":"cus_001",` +
		`"total_amount":3000,"occurred_at":"2023-11-14T22:13:20Z"}`

	var e events.OrderPlaced
	if err := json.Unmarshal([]byte(payload), &e); err != nil {
		t.Fatal(err)
	}
	want := events.OrderPlaced{
		OrderID:     "ord_001",
		CustomerID:  "cus_001",
		TotalAmount: 3000,
		OccurredAt:  time.Unix(1700000000, 0).UTC(),
	}
	if e != want {
		t.Errorf("got %+v, want %+v", e, want)
	}
}
