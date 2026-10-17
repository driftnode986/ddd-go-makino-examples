package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/infrastructure"
)

type fixedClock struct{ at time.Time }

func (c fixedClock) Now() time.Time { return c.at }

func newPlaceOrderUseCase(repo *infrastructure.MemoryOrderRepository) *application.PlaceOrderUseCase {
	return &application.PlaceOrderUseCase{
		OrderRepo: repo,
		Clock:     fixedClock{at: time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC)},
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func appCode(t *testing.T, err error) string {
	t.Helper()
	var appErr *application.ApplicationError
	if !errors.As(err, &appErr) {
		t.Fatalf("err = %v, want *ApplicationError", err)
	}
	return appErr.Code
}

func TestPlaceOrderUseCase_SavesOrder(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryOrderRepository()
	uc := newPlaceOrderUseCase(repo)

	res, err := uc.Execute(ctx, application.PlaceOrderCommand{
		CustomerID: "cus_001",
		Items: []application.OrderItemInput{
			{ProductID: "prod_001", Quantity: 2, UnitPrice: 1500},
		},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.Total != 3000 {
		t.Errorf("total = %d, want 3000", res.Total)
	}

	get := &application.GetOrderQueryHandler{OrderRepo: repo}
	summary, err := get.Execute(ctx, application.GetOrderQuery{OrderID: res.OrderID})
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if summary.Status != "draft" || summary.Total != 3000 {
		t.Errorf("summary = %+v", summary)
	}

	list := &application.ListOrdersByCustomerHandler{OrderRepo: repo}
	orders, err := list.Execute(ctx, application.ListOrdersByCustomerQuery{CustomerID: "cus_001"})
	if err != nil {
		t.Fatalf("ListOrders: %v", err)
	}
	if len(orders) != 1 || orders[0].OrderID != res.OrderID {
		t.Errorf("orders = %+v", orders)
	}
}

func TestPlaceOrderUseCase_InvalidInput(t *testing.T) {
	uc := newPlaceOrderUseCase(infrastructure.NewMemoryOrderRepository())
	cases := map[string]application.PlaceOrderCommand{
		"empty customer": {Items: []application.OrderItemInput{{ProductID: "p", Quantity: 1, UnitPrice: 1}}},
		"no items":       {CustomerID: "cus_001"},
		"zero quantity": {CustomerID: "cus_001", Items: []application.OrderItemInput{
			{ProductID: "p", Quantity: 0, UnitPrice: 1},
		}},
	}
	for name, cmd := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), cmd)
			if got := appCode(t, err); got != "INVALID_INPUT" {
				t.Errorf("code = %q, want INVALID_INPUT", got)
			}
		})
	}
}

func TestGetOrderQueryHandler_NotFound(t *testing.T) {
	get := &application.GetOrderQueryHandler{OrderRepo: infrastructure.NewMemoryOrderRepository()}
	_, err := get.Execute(context.Background(), application.GetOrderQuery{OrderID: "ord_missing"})
	if got := appCode(t, err); got != "NOT_FOUND" {
		t.Errorf("code = %q, want NOT_FOUND", got)
	}
}
