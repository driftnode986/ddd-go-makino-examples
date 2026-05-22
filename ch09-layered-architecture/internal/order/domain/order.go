package domain

import (
	"slices"
	"time"
)

type OrderStatus string

const (
	StatusDraft     OrderStatus = "draft"
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	productID ProductID
	quantity  int
	unitPrice Money
}

func NewOrderItem(productID ProductID, quantity int, unitPrice Money) (OrderItem, error) {
	if quantity <= 0 {
		return OrderItem{}, ErrInvalidQuantity
	}
	return OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
	}, nil
}

func (i OrderItem) ProductID() ProductID { return i.productID }
func (i OrderItem) Quantity() int        { return i.quantity }
func (i OrderItem) UnitPrice() Money     { return i.unitPrice }
func (i OrderItem) Subtotal() Money      { return i.unitPrice.MulInt(i.quantity) }

type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     OrderStatus
	placedAt   time.Time
}

func NewOrder(customerID CustomerID, placedAt time.Time) *Order {
	return &Order{
		id:         NewOrderID(),
		customerID: customerID,
		status:     StatusDraft,
		placedAt:   placedAt,
	}
}

func (o *Order) ID() OrderID            { return o.id }
func (o *Order) CustomerID() CustomerID { return o.customerID }
func (o *Order) Status() OrderStatus    { return o.status }
func (o *Order) PlacedAt() time.Time    { return o.placedAt }
func (o *Order) Items() []OrderItem     { return slices.Clone(o.items) }

func (o *Order) AddItem(productID ProductID, quantity int, unitPrice Money) error {
	if o.status != StatusDraft {
		return ErrCannotModify
	}
	item, err := NewOrderItem(productID, quantity, unitPrice)
	if err != nil {
		return err
	}
	o.items = append(o.items, item)
	return nil
}

func (o *Order) Confirm() error {
	if o.status != StatusDraft {
		return ErrCannotModify
	}
	o.status = StatusConfirmed
	return nil
}

func (o *Order) Total() Money {
	if len(o.items) == 0 {
		return MustMoney(0, "JPY")
	}
	total := MustMoney(0, o.items[0].UnitPrice().Currency())
	for _, item := range o.items {
		total, _ = total.Add(item.Subtotal())
	}
	return total
}
