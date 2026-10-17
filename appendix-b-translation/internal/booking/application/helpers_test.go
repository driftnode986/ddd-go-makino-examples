package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

// newPlaceOrderCommand は 1500 円の商品 2 個の注文。
func newPlaceOrderCommand() application.PlaceOrderCommand {
	return application.PlaceOrderCommand{
		CustomerID: "cus_001",
		Items: []domain.OrderItem{
			{ProductID: "prod_001", Quantity: 2, UnitPrice: 1500},
		},
	}
}

// buildExpectedOrder は Execute が保存するはずの注文を組み立てる。
func buildExpectedOrder(
	id domain.OrderID,
	cmd application.PlaceOrderCommand,
	placedAt time.Time,
) *domain.Order {
	return domain.ReconstructOrder(id, cmd.CustomerID, cmd.Items,
		domain.OrderStatusPlaced, placedAt)
}

func TestPlaceOrderUseCase_ReturnsTotalAndErrors(t *testing.T) {
	ctx := context.Background()
	clock := test.NewFixedClock(time.Now())
	uc := application.NewPlaceOrderUseCase(
		infra.NewInMemoryOrderRepository(),
		domain.NewPricingService(10),
		clock,
	)

	res, err := uc.Execute(ctx, newPlaceOrderCommand())
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 3300 {
		t.Errorf("total = %d, want 3300", res.Total)
	}
	_, err = uc.Execute(
		ctx,
		application.PlaceOrderCommand{CustomerID: "cus_001"},
	)
	if !errors.Is(err, domain.ErrInvalidOrder) {
		t.Errorf("no items: err = %v, want ErrInvalidOrder", err)
	}
}
