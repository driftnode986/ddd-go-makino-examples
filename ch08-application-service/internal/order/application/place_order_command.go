package application

// PlaceOrderCommand はユースケースの入力 DTO。
// Presentation 層が HTTP リクエストや gRPC リクエストから組み立てて渡す。
type PlaceOrderCommand struct {
	CustomerID string
	Items      []OrderItemRequest
}

type OrderItemRequest struct {
	ProductID string
	Quantity  int
}

// PlaceOrderResult はユースケースの出力 DTO。
type PlaceOrderResult struct {
	OrderID string
	Total   int64
}
