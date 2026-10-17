package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type OrderID string

func (id OrderID) String() string { return string(id) }

type CustomerID string

// OrderItem は注文明細。値として扱うので公開フィールドで持つ。
type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice int64
}

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "placed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order は注文の集約。
type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     OrderStatus
	placedAt   time.Time
}

// PlaceOrder は明細と顧客から新しい注文を作る。
func PlaceOrder(
	items []OrderItem,
	customerID CustomerID,
	now time.Time,
) (*Order, error) {
	if customerID == "" {
		return nil, fmt.Errorf(
			"%w: customer id required",
			ErrInvalidOrder,
		)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: items required", ErrInvalidOrder)
	}
	for _, it := range items {
		if it.Quantity <= 0 || it.UnitPrice < 0 {
			return nil, fmt.Errorf(
				"%w: invalid item %q",
				ErrInvalidOrder,
				it.ProductID,
			)
		}
	}
	return &Order{
		id:         newOrderID(),
		customerID: customerID,
		items:      append([]OrderItem(nil), items...),
		status:     OrderStatusPlaced,
		placedAt:   now,
	}, nil
}

// ReconstructOrder は保存済みの値から注文を組み立て直す（リポジトリとテストが使う）。
func ReconstructOrder(
	id OrderID,
	customerID CustomerID,
	items []OrderItem,
	status OrderStatus,
	placedAt time.Time,
) *Order {
	return &Order{id: id, customerID: customerID,
		items: append(
			[]OrderItem(nil),
			items...), status: status, placedAt: placedAt}
}

func (o *Order) ID() OrderID            { return o.id }
func (o *Order) CustomerID() CustomerID { return o.customerID }

func (o *Order) Items() []OrderItem {
	return append([]OrderItem(nil), o.items...)
}

func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) PlacedAt() time.Time { return o.placedAt }

// Subtotal は税抜きの合計金額（円）。
func (o *Order) Subtotal() int64 {
	var sum int64
	for _, it := range o.items {
		sum += it.UnitPrice * int64(it.Quantity)
	}
	return sum
}

// Cancel は注文を取り消す。取り消し済みなら ErrInvalidStatus。
func (o *Order) Cancel() error {
	if o.status != OrderStatusPlaced {
		return ErrInvalidStatus
	}
	o.status = OrderStatusCancelled
	return nil
}

func newOrderID() OrderID {
	b := make([]byte, 8)
	_, _ = rand.Read(b) // crypto/rand.Read はエラーを返さない（Go 1.24 以降）
	return OrderID("ord_" + hex.EncodeToString(b))
}
