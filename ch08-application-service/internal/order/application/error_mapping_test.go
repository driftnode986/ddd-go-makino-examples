package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/infrastructure"
)

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
		t.Errorf(
			"Code = %q, want %q",
			appErr.Code,
			"INVALID_QUANTITY",
		)
	}
	// Unwrap で原因のドメインエラーまでたどれる
	if !errors.Is(err, domain.ErrInvalidQuantity) {
		t.Errorf(
			"errors.Is(err, ErrInvalidQuantity) = false: %v",
			err,
		)
	}
	if got, want := appErr.Error(),
		"INVALID_QUANTITY: 数量が正しくありません"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
