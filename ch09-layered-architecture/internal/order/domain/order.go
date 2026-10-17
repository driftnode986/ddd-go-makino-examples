// #region order
package domain

import (
	"errors"
	"time"
)

type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     Status
	createdAt  time.Time
}

func NewOrder(custID CustomerID, now time.Time) (*Order, error) {
	if custID.IsZero() {
		return nil, errors.New("customer id required")
	}
	return &Order{
		id:         NewOrderID(),
		customerID: custID,
		status:     StatusDraft,
		createdAt:  now,
	}, nil
}

// #endregion

func (o *Order) ID() OrderID            { return o.id }
func (o *Order) CustomerID() CustomerID { return o.customerID }
func (o *Order) Status() Status         { return o.status }
func (o *Order) CreatedAt() time.Time   { return o.createdAt }

// Items は明細のコピーを返す（内部スライスを外から書き換えさせない）。
func (o *Order) Items() []OrderItem {
	return append([]OrderItem(nil), o.items...)
}

func (o *Order) AddItem(
	productID ProductID,
	quantity int,
	unitPrice Money,
) error {
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
