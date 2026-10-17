// #region test-imports
package application_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/events"
)

// #endregion

// #region idempotent-test
func TestHandleOrderPlaced_Idempotent(t *testing.T) {
	mailer := &stubMailer{}
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
	if err := handler.HandleWithIdempotency(ctx, 42, e); err != nil {
		t.Fatal(err)
	}

	if mailer.sendCount != 1 {
		t.Errorf("mailer called %d times, want 1", mailer.sendCount)
	}
}

// #endregion
