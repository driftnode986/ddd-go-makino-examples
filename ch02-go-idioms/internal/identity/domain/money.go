package domain

import "fmt"

type Currency string

type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("amount must be non-negative")
	}
	return Money{amount: amount, currency: currency}, nil
}

// Add は新しい Money を返す (状態変更ではない)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch")
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Equals は値比較のみ
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
