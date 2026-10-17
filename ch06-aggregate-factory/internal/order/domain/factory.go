// #region new-order
package domain

import (
	"errors"
	"slices"
	"time"
)

// NewOrder はファクトリ関数。依存ゼロのケース。
func NewOrder(id OrderID, customerID CustomerID) (*Order, error) {
	if id.IsZero() || customerID.IsZero() {
		return nil, errors.New("id, customerID required")
	}
	return &Order{
		id:         id,
		customerID: customerID,
		items:      nil,
		status:     StatusDraft,
	}, nil
}

// #endregion

// #region order-factory
// IDGenerator は OrderID を新規生成する依存。
type IDGenerator interface {
	NewOrderID() OrderID
}

// OrderFactory はファクトリ構造体。依存 (ID 生成、Clock) を持つ。
type OrderFactory struct {
	idGen IDGenerator
	clock Clock
}

// NewOrderFactory はファクトリ構造体を生成する。
func NewOrderFactory(idGen IDGenerator, clock Clock) *OrderFactory {
	return &OrderFactory{idGen: idGen, clock: clock}
}

// Create は ID 生成と検証を一括で行い、Order を返す。
func (f *OrderFactory) Create(customerID CustomerID) (*Order, error) {
	if customerID.IsZero() {
		return nil, errors.New("customerID required")
	}
	return &Order{
		id:         f.idGen.NewOrderID(),
		customerID: customerID,
		items:      nil,
		status:     StatusDraft,
	}, nil
}

// #endregion

// #region order-options
// OrderOption は Order 生成時の追加設定。
type OrderOption func(*Order)

// WithNotes は注文に備考を追加するオプション。
func WithNotes(notes string) OrderOption {
	return func(o *Order) { o.notes = notes }
}

// WithDiscountRate は注文に割引率を適用するオプション。
func WithDiscountRate(rate DiscountRate) OrderOption {
	return func(o *Order) { o.discountRate = rate }
}

// #endregion

// #region new-order-with-options
// NewOrderWithOptions は Option 関数を受け取って Order を生成する。
func NewOrderWithOptions(
	id OrderID, customerID CustomerID, opts ...OrderOption,
) (*Order, error) {
	o := &Order{
		id:         id,
		customerID: customerID,
		items:      nil,
		status:     StatusDraft,
	}
	for _, opt := range opts {
		opt(o)
	}
	if err := o.validate(); err != nil {
		return nil, err
	}
	return o, nil
}

// validate は集約全体の不変条件をチェックする内部メソッド。
func (o *Order) validate() error {
	if o.id.IsZero() || o.customerID.IsZero() {
		return errors.New("id, customerID required")
	}
	return nil
}

// #endregion

// #region reconstruct-order
// ReconstructOrder は永続化層から集約を復元するためのファクトリ。
// 入力検証はスキップし、DB に既に存在する状態をそのままメモリに復元する。
// リポジトリ層からのみ呼ぶ前提。
func ReconstructOrder(
	id OrderID,
	customerID CustomerID,
	items []OrderItem,
	status OrderStatus,
	placedAt time.Time,
) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		items:      slices.Clone(items),
		status:     status,
		placedAt:   placedAt,
	}
}

// #endregion
