// #region discount-rate
package domain

import (
	"errors"
	"fmt"
)

// DiscountRate は割引率を表す値オブジェクト (0.0-1.0)。
type DiscountRate struct {
	value float64
}

// ErrInvalidDiscountRate は不正な割引率を表すドメインエラー。
var ErrInvalidDiscountRate = errors.New("invalid discount rate")

// NewDiscountRate は割引率を生成する。
func NewDiscountRate(v float64) (DiscountRate, error) {
	if v < 0 || v > 1 {
		return DiscountRate{}, fmt.Errorf("%w: must be in [0, 1]", ErrInvalidDiscountRate)
	}
	return DiscountRate{value: v}, nil
}

// Value は割引率を float64 で返す。
func (d DiscountRate) Value() float64 { return d.value }

// #endregion

// #region calculate
// CalculateDiscountedPrice は割引適用後の価格を計算するドメインサービス。
// 依存ゼロの純粋計算なので関数で書く。
func CalculateDiscountedPrice(original Money, rate DiscountRate) (Money, error) {
	if rate.value < 0 || rate.value > 1 {
		return Money{}, ErrInvalidDiscountRate
	}
	// 整数演算で精度を保つ: amount * (1000 - rate*1000) / 1000
	rateInThousandths := int64(rate.value * 1000)
	discounted := original.amount * (1000 - rateInThousandths) / 1000
	return NewMoney(discounted, original.currency)
}

// #endregion
