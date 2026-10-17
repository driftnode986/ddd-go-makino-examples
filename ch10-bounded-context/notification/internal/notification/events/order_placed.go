package events

import "time"

// OrderPlaced は Booking BC が発火する
// 「公表された言語」としての DTO。
// Booking BC の domain.OrderPlaced をそのまま import せず、
// Notification BC 側で独立した型を持つ。
type OrderPlaced struct {
	OrderID     string    `json:"order_id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount int64     `json:"total_amount"`
	OccurredAt  time.Time `json:"occurred_at"`
}
