package infrastructure

import (
	"context"
	"errors"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

// InMemoryOrderRepository はテストで使う in-memory のリポジトリ。
type InMemoryOrderRepository struct {
	mu     sync.Mutex
	orders map[domain.OrderID]*domain.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: map[domain.OrderID]*domain.Order{},
	}
}

// ErrDuplicateOrder は同じ ID の注文を 2 回 Add したときのエラー。
var ErrDuplicateOrder = errors.New("order already exists")

func (r *InMemoryOrderRepository) Add(
	_ context.Context,
	o *domain.Order,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID()]; ok {
		return ErrDuplicateOrder
	}
	r.orders[o.ID()] = o
	return nil
}

func (r *InMemoryOrderRepository) FindByID(
	_ context.Context,
	id domain.OrderID,
) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return o, nil
}

var _ domain.OrderRepository = (*InMemoryOrderRepository)(nil)
