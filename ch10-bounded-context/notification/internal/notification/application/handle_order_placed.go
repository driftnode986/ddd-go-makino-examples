// #region handler-types
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/events"
)

type Mailer interface {
	Send(ctx context.Context, customerID, body string) error
}

type OrderPlacedHandler struct {
	Mailer      Mailer
	Idempotency IdempotencyStore
	Logger      *slog.Logger
}

// #endregion

// #region handle
func (h *OrderPlacedHandler) Handle(
	ctx context.Context, e events.OrderPlaced,
) error {
	body := fmt.Sprintf("ご注文 %s を承りました (合計 %d 円)",
		e.OrderID, e.TotalAmount)
	if err := h.Mailer.Send(ctx, e.CustomerID, body); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	h.Logger.InfoContext(ctx, "order placed notification sent",
		slog.String("customer_id", e.CustomerID),
		slog.String("order_id", e.OrderID),
	)
	return nil
}

// #endregion

// HandleWithIdempotency は冪等性ストアで重複処理を排除する。
// #region handle-with-idempotency
func (h *OrderPlacedHandler) HandleWithIdempotency(
	ctx context.Context, eventID int64, e events.OrderPlaced,
) error {
	processed, err := h.Idempotency.AlreadyProcessed(ctx, eventID)
	if err != nil {
		return fmt.Errorf("idempotency check: %w", err)
	}
	if processed {
		h.Logger.InfoContext(ctx, "event already processed, skip",
			slog.Int64("event_id", eventID))
		return nil
	}

	if err := h.Handle(ctx, e); err != nil {
		return err
	}

	return h.Idempotency.MarkProcessed(ctx, eventID)
}

// #endregion
