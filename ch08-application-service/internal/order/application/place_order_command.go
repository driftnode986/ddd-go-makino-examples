package application

type PlaceOrderCommand struct {
	CustomerID string
	Items      []OrderItemRequest
}

type OrderItemRequest struct {
	ProductID string
	Quantity  int
}

type PlaceOrderResult struct {
	OrderID string
	Total   int64
}
