// #region order-item-a
package domain

import (
	"errors"
)

// OrderItem は注文明細行 (集約内部の値オブジェクト)。
type OrderItem struct {
	productID ProductID
	quantity  int
	unitPrice Money
}

// ErrInvalidUnitPrice は不正な単価のドメインエラー。
var ErrInvalidUnitPrice = errors.New("unit price must be positive")

// #endregion
// #region order-item-b
// NewOrderItem は OrderItem を生成する。
func NewOrderItem(
	productID ProductID,
	quantity int,
	unitPrice Money,
) (OrderItem, error) {
	if productID.IsZero() {
		return OrderItem{}, errors.New("productID required")
	}
	if quantity <= 0 {
		return OrderItem{}, ErrInvalidQuantity
	}
	if unitPrice.Amount() <= 0 {
		return OrderItem{}, ErrInvalidUnitPrice
	}
	return OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
	}, nil
}

// #endregion

// #region order-item-methods
// ProductID は商品 ID。
func (i OrderItem) ProductID() ProductID { return i.productID }

// Quantity は数量。
func (i OrderItem) Quantity() int { return i.quantity }

// UnitPrice は単価。
func (i OrderItem) UnitPrice() Money { return i.unitPrice }

// Subtotal は小計 (unitPrice * quantity)。
func (i OrderItem) Subtotal() Money {
	return Money{
		amount:   i.unitPrice.Amount() * int64(i.quantity),
		currency: i.unitPrice.Currency(),
	}
}

// Equals は同値比較 (値オブジェクトとしての等価性)。
func (i OrderItem) Equals(other OrderItem) bool {
	return i.productID.Equals(other.productID) &&
		i.quantity == other.quantity &&
		i.unitPrice.Equals(other.unitPrice)
}

// #endregion
