package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure/sqlcgen"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/pgxtx"
)

// OrderRepositoryPG は orders テーブルに注文を保存する。明細は JSONB の 1 列に入れる。
type OrderRepositoryPG struct {
	queries *sqlcgen.Queries
}

func NewOrderRepositoryPG(db *pgxpool.Pool) *OrderRepositoryPG {
	return &OrderRepositoryPG{queries: sqlcgen.New(db)}
}

func (r *OrderRepositoryPG) Add(ctx context.Context, o *domain.Order) error {
	items, err := json.Marshal(o.Items())
	if err != nil {
		return fmt.Errorf("marshal items: %w", err)
	}
	return r.q(ctx).InsertOrder(ctx, sqlcgen.InsertOrderParams{
		ID:         o.ID().String(),
		CustomerID: string(o.CustomerID()),
		Status:     string(o.Status()),
		Items:      items,
		PlacedAt:   pgtype.Timestamptz{Time: o.PlacedAt(), Valid: true},
	})
}

func (r *OrderRepositoryPG) FindByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	row, err := r.q(ctx).GetOrderByID(ctx, id.String())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	var items []domain.OrderItem
	if err := json.Unmarshal(row.Items, &items); err != nil {
		return nil, fmt.Errorf("unmarshal items: %w", err)
	}
	return domain.ReconstructOrder(domain.OrderID(row.ID), domain.CustomerID(row.CustomerID),
		items, domain.OrderStatus(row.Status), row.PlacedAt.Time), nil
}

func (r *OrderRepositoryPG) q(ctx context.Context) *sqlcgen.Queries {
	if tx, ok := pgxtx.TxFromContext(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

var _ domain.OrderRepository = (*OrderRepositoryPG)(nil)
