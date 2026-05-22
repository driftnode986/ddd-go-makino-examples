package domain

import (
	"errors"
	"slices"
	"time"
)

// OrderStatus は注文の状態。
type OrderStatus string

const (
	StatusDraft     OrderStatus = "draft"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusCancelled OrderStatus = "cancelled"
)

// Order は注文集約のルート。
type Order struct {
	id           OrderID
	customerID   CustomerID
	items        []OrderItem
	status       OrderStatus
	placedAt     time.Time
	notes        string
	discountRate DiscountRate
}

// ドメインエラー。
var (
	ErrCannotModifyConfirmed = errors.New("cannot modify confirmed order")
	ErrInvalidQuantity       = errors.New("quantity must be positive")
	ErrEmptyOrder            = errors.New("cannot confirm empty order")
	ErrItemNotFound          = errors.New("item not found in order")
)

// AddItem は明細行を追加する。draft 状態でのみ可能。
func (o *Order) AddItem(productID ProductID, quantity int, unitPrice Money) error {
	if o.status != StatusDraft {
		return ErrCannotModifyConfirmed
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	item, err := NewOrderItem(productID, quantity, unitPrice)
	if err != nil {
		return err
	}
	o.items = append(o.items, item)
	return nil
}

// RemoveItem は明細行を削除する。draft 状態でのみ可能。
func (o *Order) RemoveItem(productID ProductID) error {
	if o.status != StatusDraft {
		return ErrCannotModifyConfirmed
	}
	for i, item := range o.items {
		if item.ProductID().Equals(productID) {
			o.items = slices.Delete(o.items, i, i+1)
			return nil
		}
	}
	return ErrItemNotFound
}

// Confirm は注文を確定する。
func (o *Order) Confirm(clock Clock) error {
	if o.status != StatusDraft {
		return ErrCannotModifyConfirmed
	}
	if len(o.items) == 0 {
		return ErrEmptyOrder
	}
	o.status = StatusConfirmed
	o.placedAt = clock.Now()
	return nil
}

// Items は明細行の読み取り専用ビュー (defensive copy)。
func (o Order) Items() []OrderItem {
	return slices.Clone(o.items)
}

// 状態取得。

// ID は OrderID を返す。
func (o Order) ID() OrderID { return o.id }

// CustomerID は CustomerID を返す。
func (o Order) CustomerID() CustomerID { return o.customerID }

// Status は OrderStatus を返す。
func (o Order) Status() OrderStatus { return o.status }

// PlacedAt は確定日時を返す。
func (o Order) PlacedAt() time.Time { return o.placedAt }

// Notes は備考を返す。
func (o Order) Notes() string { return o.notes }

// DiscountRate は割引率を返す。
func (o Order) DiscountRate() DiscountRate { return o.discountRate }

// TotalAmount は明細の合計金額。
func (o Order) TotalAmount() (Money, error) {
	if len(o.items) == 0 {
		return Money{}, nil
	}
	total := Money{currency: o.items[0].Subtotal().Currency()}
	for _, item := range o.items {
		added, err := total.Add(item.Subtotal())
		if err != nil {
			return Money{}, err
		}
		total = added
	}
	return total, nil
}
