// #region place-order-usecase
package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

type PlaceOrderUseCase struct {
	orders  domain.OrderRepository
	pricing domain.PricingService
	clock   Clock
}

func NewPlaceOrderUseCase(orders domain.OrderRepository, pricing domain.PricingService, clock Clock) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{orders: orders, pricing: pricing, clock: clock}
}

// #endregion

// Execute は注文を作って保存し、税込みの合計を返す。
func (uc *PlaceOrderUseCase) Execute(ctx context.Context, cmd PlaceOrderCommand) (PlaceOrderResult, error) {
	order, err := domain.PlaceOrder(cmd.Items, cmd.CustomerID, uc.clock.Now())
	if err != nil {
		return PlaceOrderResult{}, err
	}
	if err := uc.orders.Add(ctx, order); err != nil {
		return PlaceOrderResult{}, err
	}
	return PlaceOrderResult{OrderID: order.ID(), Total: uc.pricing.Total(order)}, nil
}
