package domain

import "errors"

var (
	ErrNotFound        = errors.New("order not found")
	ErrDuplicate       = errors.New("order already exists")
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrCannotModify    = errors.New("order cannot be modified in current state")
)
