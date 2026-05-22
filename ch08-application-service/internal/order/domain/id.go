package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidID = errors.New("invalid id")

type CustomerID struct {
	value string
}

func NewCustomerID(s string) (CustomerID, error) {
	if s == "" {
		return CustomerID{}, fmt.Errorf("customer id: %w", ErrInvalidID)
	}
	return CustomerID{value: s}, nil
}

func (id CustomerID) String() string { return id.value }

type ProductID struct {
	value string
}

func NewProductID(s string) (ProductID, error) {
	if s == "" {
		return ProductID{}, fmt.Errorf("product id: %w", ErrInvalidID)
	}
	return ProductID{value: s}, nil
}

func (id ProductID) String() string { return id.value }

type OrderID struct {
	value string
}

func NewOrderID(s string) (OrderID, error) {
	if s == "" {
		return OrderID{}, fmt.Errorf("order id: %w", ErrInvalidID)
	}
	return OrderID{value: s}, nil
}

func (id OrderID) String() string { return id.value }
