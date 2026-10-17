// #region usecase
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

// PlaceOrderUseCase は注文を確定するユースケース。
type PlaceOrderUseCase struct {
	orders  domain.Repository
	pricing PricingPort
	clock   Clock
	logger  *slog.Logger
}

func NewPlaceOrderUseCase(
	orders domain.Repository,
	pricing PricingPort,
	clock Clock,
	logger *slog.Logger,
) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		orders: orders, pricing: pricing, clock: clock, logger: logger,
	}
}

// #endregion

// #region execute
// Execute は注文を確定する。1 ユースケース = 1 メソッド。
func (uc *PlaceOrderUseCase) Execute(
	ctx context.Context, cmd PlaceOrderCommand,
) (PlaceOrderResult, error) {
	customerID, err := domain.NewCustomerID(cmd.CustomerID)
	if err != nil {
		return PlaceOrderResult{}, fmt.Errorf("invalid customer id: %w", err)
	}

	o := domain.NewOrder(customerID, uc.clock.Now())
	for _, item := range cmd.Items {
		productID, err := domain.NewProductID(item.ProductID)
		if err != nil {
			return PlaceOrderResult{}, fmt.Errorf("invalid product id: %w", err)
		}
		unitPrice, err := uc.pricing.UnitPrice(ctx, productID)
		if err != nil {
			return PlaceOrderResult{}, fmt.Errorf("fetch unit price: %w", err)
		}
		if err := o.AddItem(productID, item.Quantity, unitPrice); err != nil {
			return PlaceOrderResult{}, mapDomainError(err)
		}
	}

	// #endregion

	// #region execute-save
	if err := uc.orders.Add(ctx, o); err != nil {
		return PlaceOrderResult{}, fmt.Errorf("save order: %w", err)
	}

	uc.logger.InfoContext(ctx, "order placed",
		slog.String("order_id", o.ID().String()),
		slog.String("customer_id", customerID.String()),
		slog.Int("item_count", len(cmd.Items)),
		slog.Int64("total", o.Total().Amount()),
	)

	// #region result
	return PlaceOrderResult{
		OrderID: o.ID().String(),
		Total:   o.Total().Amount(),
	}, nil
	// #endregion
}

// #endregion
