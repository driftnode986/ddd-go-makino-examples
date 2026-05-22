package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
)

// OrderItemInput はユースケースが受け取る明細の入力。
// Presentation 層の DTO とは別の型として宣言する。
type OrderItemInput struct {
	ProductID string
	Quantity  int
	UnitPrice int64
}

type PlaceOrderCommand struct {
	CustomerID string
	Items      []OrderItemInput
}

type PlaceOrderResult struct {
	OrderID string
	Total   int64
}

type PlaceOrderUseCase struct {
	OrderRepo domain.Repository
	Clock     Clock
	Logger    *slog.Logger
}

func (uc *PlaceOrderUseCase) Execute(
	ctx context.Context, cmd PlaceOrderCommand,
) (PlaceOrderResult, error) {
	customerID, err := domain.NewCustomerID(cmd.CustomerID)
	if err != nil {
		return PlaceOrderResult{}, &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "customer id is invalid",
			Cause:   err,
		}
	}

	if len(cmd.Items) == 0 {
		return PlaceOrderResult{}, &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "items must not be empty",
		}
	}

	o := domain.NewOrder(customerID, uc.Clock.Now())
	for _, item := range cmd.Items {
		productID, err := domain.NewProductID(item.ProductID)
		if err != nil {
			return PlaceOrderResult{}, &ApplicationError{
				Code:    "INVALID_INPUT",
				Message: "product id is invalid",
				Cause:   err,
			}
		}
		unitPrice, err := domain.NewMoney(item.UnitPrice, "JPY")
		if err != nil {
			return PlaceOrderResult{}, fmt.Errorf("unit price: %w", err)
		}
		if err := o.AddItem(productID, item.Quantity, unitPrice); err != nil {
			return PlaceOrderResult{}, mapDomainError(err)
		}
	}

	if err := uc.OrderRepo.Add(ctx, o); err != nil {
		return PlaceOrderResult{}, fmt.Errorf("save order: %w", err)
	}

	uc.Logger.InfoContext(ctx, "order placed",
		slog.String("order_id", o.ID().String()),
		slog.String("customer_id", customerID.String()),
		slog.Int("item_count", len(cmd.Items)),
		slog.Int64("total", o.Total().Amount()),
	)

	return PlaceOrderResult{
		OrderID: o.ID().String(),
		Total:   o.Total().Amount(),
	}, nil
}
