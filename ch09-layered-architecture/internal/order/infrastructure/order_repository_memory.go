package infrastructure

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
)

type MemoryOrderRepository struct {
	mu     sync.Mutex
	orders map[domain.OrderID]*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{orders: map[domain.OrderID]*domain.Order{}}
}

func (r *MemoryOrderRepository) Add(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[o.ID()]; exists {
		return domain.ErrDuplicate
	}
	r.orders[o.ID()] = o
	return nil
}

func (r *MemoryOrderRepository) Save(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID()]; !ok {
		return domain.ErrNotFound
	}
	r.orders[o.ID()] = o
	return nil
}

func (r *MemoryOrderRepository) FindByID(
	ctx context.Context, id domain.OrderID,
) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (r *MemoryOrderRepository) ListByCustomer(
	ctx context.Context, customerID domain.CustomerID,
) ([]*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*domain.Order, 0)
	for _, o := range r.orders {
		if o.CustomerID() == customerID {
			result = append(result, o)
		}
	}
	return result, nil
}

var _ domain.Repository = (*MemoryOrderRepository)(nil)
