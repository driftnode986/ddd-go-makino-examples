// Package money demonstrates Go 1.26 self-referencing generic constraints.
package money

import "fmt"

// Adder is a constraint that requires the type to be addable to itself.
// Self-reference in type parameter lists is allowed in Go 1.26.
type Adder[A any] interface {
	Add(A) A
}

// Currency is a phantom type marker for compile-time currency safety.
type Currency interface {
	currencyCode() string
}

type JPY struct{}

func (JPY) currencyCode() string { return "JPY" }

type USD struct{}

func (USD) currencyCode() string { return "USD" }

// Money carries a currency phantom type C so the compiler rejects cross-currency addition.
type Money[C Currency] struct {
	amount int64
}

func NewMoney[C Currency](amount int64) Money[C] {
	return Money[C]{amount: amount}
}

// Add only accepts the same currency type parameter at compile time.
func (m Money[C]) Add(other Money[C]) Money[C] {
	return Money[C]{amount: m.amount + other.amount}
}

func (m Money[C]) Amount() int64 {
	return m.amount
}

func (m Money[C]) String() string {
	var c C
	return fmt.Sprintf("%d %s", m.amount, c.currencyCode())
}

// Compile-time check that Money[JPY] satisfies the self-referencing Adder constraint.
var _ Adder[Money[JPY]] = Money[JPY]{}
