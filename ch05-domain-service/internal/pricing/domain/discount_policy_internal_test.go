package domain

import (
	"errors"
	"testing"
)

// NewDiscountRate を通らずに作った範囲外の割引率も、計算の側で拒む
func TestCalculateDiscountedPrice_RejectsOutOfRangeRate(t *testing.T) {
	original, err := NewMoney(10000, JPY)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []float64{-0.1, 1.5} {
		_, err := CalculateDiscountedPrice(original, DiscountRate{value: v})
		if !errors.Is(err, ErrInvalidDiscountRate) {
			t.Errorf("rate %v: want ErrInvalidDiscountRate, got %v", v, err)
		}
	}
}
