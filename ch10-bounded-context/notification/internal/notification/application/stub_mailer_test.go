package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/events"
)

// stubMailer は Mailer の代わりに送信回数と本文を数えるテスト用の偽物。
type stubMailer struct {
	sendCount int
	lastBody  string
	err       error
}

func (m *stubMailer) Send(_ context.Context, _, body string) error {
	if m.err != nil {
		return m.err
	}
	m.sendCount++
	m.lastBody = body
	return nil
}

func newHandler(m *stubMailer) *application.OrderPlacedHandler {
	return &application.OrderPlacedHandler{
		Mailer:      m,
		Idempotency: application.NewMemoryIdempotencyStore(),
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestHandleOrderPlaced_SendsMail(t *testing.T) {
	mailer := &stubMailer{}
	e := events.OrderPlaced{
		OrderID:     "ord_001",
		CustomerID:  "cus_001",
		TotalAmount: 3000,
	}
	if err := newHandler(mailer).Handle(
		context.Background(),
		e,
	); err != nil {
		t.Fatal(err)
	}
	if mailer.sendCount != 1 ||
		!strings.Contains(mailer.lastBody, "ord_001") {
		t.Errorf(
			"sendCount = %d, body = %q",
			mailer.sendCount,
			mailer.lastBody,
		)
	}
}

func TestHandleWithIdempotency_DoesNotMarkWhenSendFails(
	t *testing.T,
) {
	mailer := &stubMailer{err: errors.New("smtp down")}
	h := newHandler(mailer)
	ctx := context.Background()
	e := events.OrderPlaced{
		OrderID:    "ord_001",
		CustomerID: "cus_001",
	}

	if err := h.HandleWithIdempotency(ctx, 7, e); err == nil {
		t.Fatal("want error when mail fails, got nil")
	}
	processed, err := h.Idempotency.AlreadyProcessed(ctx, 7)
	if err != nil {
		t.Fatal(err)
	}
	if processed {
		t.Error(
			"event must not be marked processed when sending failed",
		)
	}

	mailer.err = nil
	if err := h.HandleWithIdempotency(ctx, 7, e); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if mailer.sendCount != 1 {
		t.Errorf(
			"sendCount after retry = %d, want 1",
			mailer.sendCount,
		)
	}
}
