package domain

import (
	"errors"
	"fmt"
)

// CouponCode はクーポンコードの値オブジェクト。
type CouponCode struct {
	value string
}

// ErrInvalidCouponCode は不正なコードのドメインエラー。
var ErrInvalidCouponCode = errors.New("invalid coupon code")

// NewCouponCode はクーポンコードを生成する。
func NewCouponCode(s string) (CouponCode, error) {
	if len(s) != 8 {
		return CouponCode{}, fmt.Errorf(
			"%w: must be 8 characters",
			ErrInvalidCouponCode,
		)
	}
	return CouponCode{value: s}, nil
}

// String はコード文字列を返す。
func (c CouponCode) String() string { return c.value }

// Equals は同値比較。
func (c CouponCode) Equals(
	other CouponCode,
) bool {
	return c.value == other.value
}
