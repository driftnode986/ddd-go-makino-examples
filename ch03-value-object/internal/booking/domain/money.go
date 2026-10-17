// #region money-type
package domain

import (
	"errors"
	"fmt"
)

// Currency は通貨コード (ISO 4217)。
type Currency string

const (
	JPY Currency = "JPY"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

// Money は値オブジェクト。Amount と Currency を 1 つの概念として扱う。
// Amount は最小単位 (JPY なら円、USD ならセント) で保持し、浮動小数の誤差を避ける。
type Money struct {
	amount   int64
	currency Currency
}

// ErrCurrencyMismatch は異なる通貨同士の演算を表すドメインエラー。
var ErrCurrencyMismatch = errors.New("currency mismatch")

// ErrNegativeAmount は負の金額を表すドメインエラー。
var ErrNegativeAmount = errors.New("amount must not be negative")

// #endregion

// #region new-money
// NewMoney は最小単位の整数と通貨を受け取り、Money を返す。
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}
	if currency == "" {
		return Money{}, errors.New("currency required")
	}
	return Money{amount: amount, currency: currency}, nil
}

// #endregion

// #region money-methods
// Amount は最小単位の金額を返す。
func (m Money) Amount() int64 { return m.amount }

// Currency は通貨コードを返す。
func (m Money) Currency() Currency { return m.currency }

// Add は同通貨の Money を加算した新しい Money を返す。
// 受け手は変更されない (不変性)。
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf(
			"%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency,
		)
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Equals は値が完全に一致するときに真。
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String は人間可読な表現を返す。
func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.currency)
}

// #endregion
