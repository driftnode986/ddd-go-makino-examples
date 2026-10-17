// #region list-query-types
package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

type ListOrdersByCustomerQuery struct {
	CustomerID string
}

type OrderSummary struct {
	OrderID  string
	Total    int64
	Status   string
	PlacedAt int64
}

type ListOrdersByCustomerHandler struct {
	Orders domain.Repository
}

// #endregion

// #region list-query-execute
func (h *ListOrdersByCustomerHandler) Execute(
	ctx context.Context, q ListOrdersByCustomerQuery,
) ([]OrderSummary, error) {
	customerID, err := domain.NewCustomerID(q.CustomerID)
	if err != nil {
		return nil, err
	}
	orders, err := h.Orders.ListByCustomer(ctx, customerID)
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

// #endregion
