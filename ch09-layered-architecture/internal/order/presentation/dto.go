// #region request-dto
package presentation

import (
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
)

// PlaceOrderRequest は HTTP リクエストボディの JSON 構造。
// json タグで API 契約を表現する。
type PlaceOrderRequest struct {
	CustomerID string                  `json:"customer_id"`
	Items      []PlaceOrderItemRequest `json:"items"`
}

type PlaceOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}

// PlaceOrderResponse は HTTP レスポンスボディ。
type PlaceOrderResponse struct {
	OrderID string `json:"order_id"`
	Total   int64  `json:"total"`
}

// #endregion

// #region to-application-items
func toApplicationItems(
	reqs []PlaceOrderItemRequest,
) []application.OrderItemInput {
	out := make([]application.OrderItemInput, 0, len(reqs))
	for _, r := range reqs {
		out = append(out, application.OrderItemInput{
			ProductID: r.ProductID,
			Quantity:  r.Quantity,
			UnitPrice: r.UnitPrice,
		})
	}
	return out
}

// #endregion
