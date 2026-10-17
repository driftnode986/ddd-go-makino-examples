package application

// OrderItemInput は PlaceOrderCommand の明細 1 行。
type OrderItemInput struct {
	ProductID string
	Quantity  int
	UnitPrice int64
}

// PlaceOrderCommand は PlaceOrderUseCase の入力。
type PlaceOrderCommand struct {
	CustomerID string
	Items      []OrderItemInput
}

// PlaceOrderResult は PlaceOrderUseCase の出力。
type PlaceOrderResult struct {
	OrderID string
	Total   int64
}
