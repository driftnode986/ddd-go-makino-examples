package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/infrastructure"
)

func TestListOrdersByCustomerHandler_Empty(t *testing.T) {
	ctx := context.Background()
	orders := infrastructure.NewInMemoryOrderRepository()
	h := &application.ListOrdersByCustomerHandler{Orders: orders}

	got, err := h.Execute(ctx, application.ListOrdersByCustomerQuery{
		CustomerID: "cust_999",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d orders, want 0", len(got))
	}
}

func TestListOrdersByCustomerHandler_ReturnsSummaries(t *testing.T) {
	ctx := context.Background()
	orders := infrastructure.NewInMemoryOrderRepository()
	placedAt := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	uc := application.NewPlaceOrderUseCase(
		orders, fakePricing{unit: 1000}, fakeClock{now: placedAt},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	place := func(
		customerID string,
		qty int,
	) application.PlaceOrderResult {
		t.Helper()
		r, err := uc.Execute(ctx, application.PlaceOrderCommand{
			CustomerID: customerID,
			Items: []application.OrderItemRequest{
				{ProductID: "prod_001", Quantity: qty},
			},
		})
		if err != nil {
			t.Fatalf("Execute: %v", err)
		}
		return r
	}
	mine := place("cust_001", 2)
	place("cust_002", 1) // 別の顧客の注文は含まれない

	h := &application.ListOrdersByCustomerHandler{Orders: orders}
	got, err := h.Execute(
		ctx,
		application.ListOrdersByCustomerQuery{
			CustomerID: "cust_001",
		},
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d orders, want 1", len(got))
	}
	want := application.OrderSummary{
		OrderID:  mine.OrderID,
		Total:    2000,
		Status:   "confirmed",
		PlacedAt: placedAt.Unix(),
	}
	if got[0] != want {
		t.Errorf("got %+v, want %+v", got[0], want)
	}
}

func TestListOrdersByCustomerHandler_InvalidCustomerID(
	t *testing.T,
) {
	h := &application.ListOrdersByCustomerHandler{
		Orders: infrastructure.NewInMemoryOrderRepository(),
	}
	_, err := h.Execute(
		context.Background(),
		application.ListOrdersByCustomerQuery{},
	)
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("got %v, want ErrInvalidID", err)
	}
}
