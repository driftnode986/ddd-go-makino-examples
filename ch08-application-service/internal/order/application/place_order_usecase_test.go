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

type fakeClock struct{ now time.Time }

func (c fakeClock) Now() time.Time { return c.now }

type fakePricing struct{ unit int64 }

func (p fakePricing) UnitPrice(
	ctx context.Context, productID domain.ProductID,
) (domain.Money, error) {
	return domain.MustMoney(p.unit, "JPY"), nil
}

func TestPlaceOrderUseCase_HappyPath(t *testing.T) {
	ctx := context.Background()
	orders := infrastructure.NewInMemoryOrderRepository()
	pricing := fakePricing{unit: 1000}
	clock := fakeClock{now: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := application.NewPlaceOrderUseCase(orders, pricing, clock, logger)

	result, err := uc.Execute(ctx, application.PlaceOrderCommand{
		CustomerID: "cust_001",
		Items: []application.OrderItemRequest{
			{ProductID: "prod_001", Quantity: 2},
			{ProductID: "prod_002", Quantity: 3},
		},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if result.Total != 5000 {
		t.Errorf("Total = %d, want 5000", result.Total)
	}
	if result.OrderID == "" {
		t.Error("OrderID is empty")
	}
}

func TestPlaceOrderUseCase_InvalidCustomerID(t *testing.T) {
	ctx := context.Background()
	orders := infrastructure.NewInMemoryOrderRepository()
	uc := application.NewPlaceOrderUseCase(
		orders, fakePricing{unit: 1000}, fakeClock{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	_, err := uc.Execute(ctx, application.PlaceOrderCommand{
		CustomerID: "",
		Items: []application.OrderItemRequest{
			{ProductID: "prod_001", Quantity: 1},
		},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("got %v, want ErrInvalidID in chain", err)
	}
}

func TestPlaceOrderUseCase_InvalidQuantity(t *testing.T) {
	ctx := context.Background()
	orders := infrastructure.NewInMemoryOrderRepository()
	uc := application.NewPlaceOrderUseCase(
		orders, fakePricing{unit: 1000}, fakeClock{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	_, err := uc.Execute(ctx, application.PlaceOrderCommand{
		CustomerID: "cust_001",
		Items: []application.OrderItemRequest{
			{ProductID: "prod_001", Quantity: 0},
		},
	})
	var appErr *application.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *ApplicationError, got %T (%v)", err, err)
	}
	if appErr.Code != "INVALID_QUANTITY" {
		t.Errorf("Code = %q, want %q", appErr.Code, "INVALID_QUANTITY")
	}
}

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

func TestRequestIDContext(t *testing.T) {
	ctx := application.WithRequestID(context.Background(), "req_001")
	id, ok := application.RequestID(ctx)
	if !ok {
		t.Fatal("RequestID returned ok=false")
	}
	if id != "req_001" {
		t.Errorf("RequestID = %q, want %q", id, "req_001")
	}
}
