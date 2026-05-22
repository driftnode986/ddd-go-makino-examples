package domain

import "context"

type Repository interface {
	Add(ctx context.Context, o *Order) error
	Save(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id OrderID) (*Order, error)
	ListByCustomer(ctx context.Context, customerID CustomerID) ([]*Order, error)
}
