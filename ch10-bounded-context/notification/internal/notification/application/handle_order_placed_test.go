package application_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/events"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/infrastructure"
)

func TestHandleOrderPlaced_Success(t *testing.T) {
	mailer := infrastructure.NewMemoryMailer()
	idem := application.NewMemoryIdempotencyStore()
	handler := &application.OrderPlacedHandler{
		Mailer:      mailer,
		Idempotency: idem,
		Logger:      slog.Default(),
	}

	e := events.OrderPlaced{
		OrderID:     "ord_001",
		CustomerID:  "cus_001",
		TotalAmount: 3000,
		OccurredAt:  time.Unix(1700000000, 0),
	}

	if err := handler.Handle(context.Background(), e); err != nil {
		t.Fatal(err)
	}

	if mailer.SentCount() != 1 {
		t.Errorf("mailer called %d times, want 1", mailer.SentCount())
	}
}

func TestHandleOrderPlaced_Idempotent(t *testing.T) {
	mailer := infrastructure.NewMemoryMailer()
	idem := application.NewMemoryIdempotencyStore()
	handler := &application.OrderPlacedHandler{
		Mailer:      mailer,
		Idempotency: idem,
		Logger:      slog.Default(),
	}

	e := events.OrderPlaced{
		OrderID:     "ord_001",
		CustomerID:  "cus_001",
		TotalAmount: 3000,
	}
	ctx := context.Background()

	if err := handler.HandleWithIdempotency(ctx, 42, e); err != nil {
		t.Fatal(err)
	}
	// 同じ event ID で 2 回目 → skip される
	if err := handler.HandleWithIdempotency(ctx, 42, e); err != nil {
		t.Fatal(err)
	}

	if mailer.SentCount() != 1 {
		t.Errorf("mailer called %d times, want 1 (idempotent)", mailer.SentCount())
	}
}
