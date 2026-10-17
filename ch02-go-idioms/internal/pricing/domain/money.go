// #region money-a
package pricing

import (
	"errors"
	"fmt"
)

var (
	ErrNegativeAmount   = errors.New("amount must not be negative")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

type Currency string

type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}
	return Money{amount: amount, currency: currency}, nil
}

// #endregion
// #region money-b
// Add は新しい Money を返す (状態変更ではない)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf(
			"add: %w",
			ErrCurrencyMismatch,
		)
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract は other を引いた新しい Money を返す (負になればエラー)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf(
			"subtract: %w",
			ErrCurrencyMismatch,
		)
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

// #endregion
// #region money-c
// Amount は最小単位の金額を返す
func (m Money) Amount() int64 { return m.amount }

// Equals は値比較のみ
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// #endregion
