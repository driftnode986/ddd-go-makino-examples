package domain_test

import (
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch06-aggregate-factory/internal/order/domain"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// ExampleNewOrderWithOptions は 6.9 の呼び出し例を実行して確かめる。
func ExampleNewOrderWithOptions() {
	orderID := must(domain.NewID[domain.OrderTag]("order001"))
	customerID := must(domain.NewID[domain.CustomerTag]("cust0001"))
	rate10pct := must(domain.NewDiscountRate(0.1))

	// #region with-options-call
	order, err := domain.NewOrderWithOptions(
		orderID, customerID,
		domain.WithNotes("ギフト包装希望"),
		domain.WithDiscountRate(rate10pct),
	)
	// #endregion
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(order.Notes(), order.DiscountRate().Value())
	// Output: ギフト包装希望 0.1
}
