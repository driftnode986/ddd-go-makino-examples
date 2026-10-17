package application

import (
	"context"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

// Clock は現在時刻を返す。
type Clock interface {
	Now() time.Time
}

// TxManager は fn を 1 つのトランザクションの中で実行する。
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// PlaceOrderCommand は注文のユースケースの入力。
type PlaceOrderCommand struct {
	CustomerID domain.CustomerID
	Items      []domain.OrderItem
}

// PlaceOrderResult は PlaceOrderUseCase の出力。
type PlaceOrderResult struct {
	OrderID domain.OrderID
	Total   int64
}
