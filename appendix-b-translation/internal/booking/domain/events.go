package domain

// DomainEvent はドメインイベントが満たす interface。
type DomainEvent interface {
	EventName() string
}

// OrderPlaced は注文が確定した事実。
type OrderPlaced struct {
	OrderID OrderID
}

func (OrderPlaced) EventName() string { return "booking.OrderPlaced" }
