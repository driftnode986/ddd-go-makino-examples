package domain

import (
	"errors"
	"fmt"
)

var ErrCurrencyMismatch = errors.New("currency mismatch")

type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency required")
	}
	return Money{amount: amount, currency: currency}, nil
}

func MustMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf(
			"%w: %s vs %s",
			ErrCurrencyMismatch,
			m.currency,
			other.currency,
		)
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) MulInt(n int) Money {
	return Money{amount: m.amount * int64(n), currency: m.currency}
}
