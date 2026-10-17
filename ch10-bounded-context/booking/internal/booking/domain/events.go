// #region domain-event
package domain

import "time"

// DomainEvent はすべてのドメインイベントが満たす interface。
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() string
}

// #endregion

// #region order-placed
// OrderPlaced は注文が確定された事実を表す。
type OrderPlaced struct {
	occurredAt  time.Time
	orderID     OrderID
	customerID  CustomerID
	totalAmount int64
}

func NewOrderPlaced(o *Order, now time.Time) OrderPlaced {
	return OrderPlaced{
		occurredAt:  now,
		orderID:     o.ID(),
		customerID:  o.CustomerID(),
		totalAmount: o.Total().Amount(),
	}
}

func (e OrderPlaced) EventName() string     { return "booking.OrderPlaced" }
func (e OrderPlaced) OccurredAt() time.Time { return e.occurredAt }
func (e OrderPlaced) AggregateID() string   { return e.orderID.String() }

func (e OrderPlaced) OrderID() OrderID       { return e.orderID }
func (e OrderPlaced) CustomerID() CustomerID { return e.customerID }
func (e OrderPlaced) TotalAmount() int64     { return e.totalAmount }

// #endregion
