package domain

// OrderItem は注文明細。
type OrderItem struct {
	productID ProductID
	quantity  int
	unitPrice Money
}

func NewOrderItem(productID ProductID, quantity int, unitPrice Money) (OrderItem, error) {
	if quantity <= 0 {
		return OrderItem{}, ErrInvalidQuantity
	}
	return OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
	}, nil
}

func (i OrderItem) ProductID() ProductID { return i.productID }
func (i OrderItem) Quantity() int        { return i.quantity }
func (i OrderItem) UnitPrice() Money     { return i.unitPrice }
func (i OrderItem) Subtotal() Money      { return i.unitPrice.MulInt(i.quantity) }
