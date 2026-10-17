package domain

// DomainEvent はドメインイベントが満たす interface。
type DomainEvent interface {
	EventName() string
}

// OrderPlaced は注文が確定した事実。発生時の値を後から
// 書き換えられないよう、フィールドは非公開にする（第10章）。
type OrderPlaced struct {
	orderID OrderID
}

func NewOrderPlaced(id OrderID) OrderPlaced {
	return OrderPlaced{orderID: id}
}

func (e OrderPlaced) OrderID() OrderID { return e.orderID }

func (OrderPlaced) EventName() string {
	return "booking.OrderPlaced"
}
