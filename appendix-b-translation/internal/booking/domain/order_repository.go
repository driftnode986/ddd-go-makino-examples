package domain

import "context"

type OrderRepository interface {
	Add(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id OrderID) (*Order, error)
}
