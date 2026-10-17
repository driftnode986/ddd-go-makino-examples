// #region usecase-struct
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
)

type PlaceOrderUseCase struct {
	OrderRepo domain.Repository
	Outbox    OutboxPort
	Clock     Clock
	Logger    *slog.Logger
}

// #endregion

// #region execute-items
func (uc *PlaceOrderUseCase) Execute(
	ctx context.Context, cmd PlaceOrderCommand,
) (PlaceOrderResult, error) {
	customerID, err := domain.NewCustomerID(cmd.CustomerID)
	if err != nil {
		return PlaceOrderResult{}, &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "invalid customer id",
			Cause:   err,
		}
	}

	o := domain.NewOrder(customerID, uc.Clock.Now())
	for _, item := range cmd.Items {
		pid, err := domain.NewProductID(item.ProductID)
		if err != nil {
			return PlaceOrderResult{}, &ApplicationError{
				Code:    "INVALID_INPUT",
				Message: "invalid product id",
				Cause:   err,
			}
		}
		unitPrice, err := domain.NewMoney(item.UnitPrice, "JPY")
		if err != nil {
			return PlaceOrderResult{}, fmt.Errorf(
				"unit price: %w",
				err,
			)
		}
		if err := o.AddItem(
			pid,
			item.Quantity,
			unitPrice,
		); err != nil {
			return PlaceOrderResult{}, mapDomainError(err)
		}
	}
	// #endregion

	// #region execute-confirm-save
	if err := o.Confirm(uc.Clock.Now()); err != nil {
		return PlaceOrderResult{}, mapDomainError(err)
	}

	events := o.PullEvents()

	// リポジトリ保存と outbox 書き込みは本来 1 Tx でまとめる。
	// in-memory 実装では順次実行する。
	if err := uc.OrderRepo.Add(ctx, o); err != nil {
		return PlaceOrderResult{}, fmt.Errorf("save order: %w", err)
	}
	if err := uc.Outbox.Append(ctx, events); err != nil {
		return PlaceOrderResult{}, fmt.Errorf(
			"append outbox: %w",
			err,
		)
	}

	uc.Logger.InfoContext(ctx, "order placed",
		slog.String("order_id", o.ID().String()),
		slog.Int("event_count", len(events)),
	)

	return PlaceOrderResult{
		OrderID: o.ID().String(),
		Total:   o.Total().Amount(),
	}, nil
}

// #endregion
