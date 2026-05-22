package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
)

type ListOrdersByCustomerQuery struct {
	CustomerID string
}

type GetOrderQuery struct {
	OrderID string
}

type OrderSummary struct {
	OrderID  string
	Total    int64
	Status   string
	PlacedAt int64
}

type ListOrdersByCustomerHandler struct {
	OrderRepo domain.Repository
}

func (h *ListOrdersByCustomerHandler) Execute(
	ctx context.Context, q ListOrdersByCustomerQuery,
) ([]OrderSummary, error) {
	customerID, err := domain.NewCustomerID(q.CustomerID)
	if err != nil {
		return nil, &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "customer id is invalid",
			Cause:   err,
		}
	}
	orders, err := h.OrderRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	result := make([]OrderSummary, 0, len(orders))
	for _, o := range orders {
		result = append(result, OrderSummary{
			OrderID:  o.ID().String(),
			Total:    o.Total().Amount(),
			Status:   string(o.Status()),
			PlacedAt: o.PlacedAt().Unix(),
		})
	}
	return result, nil
}

type GetOrderQueryHandler struct {
	OrderRepo domain.Repository
}

func (h *GetOrderQueryHandler) Execute(
	ctx context.Context, q GetOrderQuery,
) (OrderSummary, error) {
	id, err := domain.ParseOrderID(q.OrderID)
	if err != nil {
		return OrderSummary{}, &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "order id is invalid",
			Cause:   err,
		}
	}
	o, err := h.OrderRepo.FindByID(ctx, id)
	if err != nil {
		return OrderSummary{}, mapDomainError(err)
	}
	return OrderSummary{
		OrderID:  o.ID().String(),
		Total:    o.Total().Amount(),
		Status:   string(o.Status()),
		PlacedAt: o.PlacedAt().Unix(),
	}, nil
}
