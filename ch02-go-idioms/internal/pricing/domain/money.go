package pricing

import "errors"

// Money は料金計算に使う金額 (最小単位の整数)。値オブジェクトとしての詳細は第3章で扱う。
type Money struct {
	amount   int64
	currency string
}

var (
	ErrNegativeAmount   = errors.New("amount must not be negative")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

// NewMoney は負でない金額の Money を返す。
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}
	return Money{amount: amount, currency: currency}, nil
}

// Amount は最小単位の金額を返す。
func (m Money) Amount() int64 { return m.amount }

// Subtract は other を引いた新しい Money を返す。結果が負になるときはエラー。
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount-other.amount, m.currency)
}
