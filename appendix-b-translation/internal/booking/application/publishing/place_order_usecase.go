// Package publishing は、注文を保存してイベントを発行する版の PlaceOrderUseCase を置く。
// application の PlaceOrderUseCase（DI の例）とは別の例なので、パッケージを分けている。
package publishing

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

type PlaceOrderCommand = application.PlaceOrderCommand

type PlaceOrderUseCase struct {
	txMgr  application.TxManager
	orders domain.OrderRepository
	events domain.EventPublisher
	clock  application.Clock
}

func NewPlaceOrderUseCase(
	txMgr application.TxManager,
	orders domain.OrderRepository,
	events domain.EventPublisher,
	clock application.Clock,
) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		txMgr:  txMgr,
		orders: orders,
		events: events,
		clock:  clock,
	}
}

// #region execute
func (uc *PlaceOrderUseCase) Execute(
	ctx context.Context,
	cmd PlaceOrderCommand,
) error {
	return uc.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		order, err := domain.PlaceOrder(
			cmd.Items,
			cmd.CustomerID,
			uc.clock.Now(),
		)
		if err != nil {
			return err
		}
		if err := uc.orders.Add(ctx, order); err != nil {
			return err
		}
		return uc.events.Publish(
			ctx,
			domain.OrderPlaced{OrderID: order.ID()},
		)
	})
}

// #endregion
