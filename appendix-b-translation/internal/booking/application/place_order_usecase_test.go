// #region saves-order-imports
package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

// #endregion

// #region saves-order
func TestPlaceOrderUseCase_SavesOrder(t *testing.T) {
	orders := infra.NewInMemoryOrderRepository()
	clock := test.NewFixedClock(
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	pricing := domain.NewPricingService(10)
	uc := application.NewPlaceOrderUseCase(orders, pricing, clock)
	ctx := context.Background()
	cmd := newPlaceOrderCommand()

	res, err := uc.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := orders.FindByID(ctx, res.OrderID)
	if err != nil {
		t.Fatalf("order not saved: %v", err)
	}
	want := buildExpectedOrder(res.OrderID, cmd, clock.Now())
	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(domain.Order{}),
	); diff != "" {
		t.Errorf("order mismatch (-want +got):\n%s", diff)
	}
}

// #endregion
