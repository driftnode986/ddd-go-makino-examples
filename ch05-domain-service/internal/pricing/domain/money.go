package domain

import (
	"errors"
	"fmt"
)

// Currency は通貨を表す値オブジェクト。
type Currency string

const (
	JPY Currency = "JPY"
	USD Currency = "USD"
)

// Money は金額の値オブジェクト。amount は最小単位 (JPY なら円、USD ならセント)。
type Money struct {
	amount   int64
	currency Currency
}

// ErrInvalidMoney は不正な金額のドメインエラー。
var ErrInvalidMoney = errors.New("invalid money")

// ErrCurrencyMismatch は通貨不一致のドメインエラー。
var ErrCurrencyMismatch = errors.New("currency mismatch")

// NewMoney は Money を生成する。
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("%w: amount must be non-negative", ErrInvalidMoney)
	}
	if currency != JPY && currency != USD {
		return Money{}, fmt.Errorf("%w: unsupported currency %q", ErrInvalidMoney, currency)
	}
	return Money{amount: amount, currency: currency}, nil
}

// Amount は最小単位での金額を返す。
func (m Money) Amount() int64 { return m.amount }

// Currency は通貨を返す。
func (m Money) Currency() Currency { return m.currency }

// Equals は同値比較。
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// Add は同じ通貨同士の加算。
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Subtract は同じ通貨同士の減算。
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	if m.amount < other.amount {
		return Money{}, fmt.Errorf("%w: result would be negative", ErrInvalidMoney)
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

// String は表示用。
func (m Money) String() string {
	return fmt.Sprintf("%d %s", m.amount, m.currency)
}
