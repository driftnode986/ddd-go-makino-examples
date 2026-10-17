package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/pricing/domain"
)

func mustMoney(
	t *testing.T,
	amount int64,
	currency domain.Currency,
) domain.Money {
	t.Helper()
	m, err := domain.NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("NewMoney(%d, %s): %v", amount, currency, err)
	}
	return m
}

func mustRate(t *testing.T, v float64) domain.DiscountRate {
	t.Helper()
	r, err := domain.NewDiscountRate(v)
	if err != nil {
		t.Fatalf("NewDiscountRate(%v): %v", v, err)
	}
	return r
}

func TestCalculateDiscountedPrice(t *testing.T) {
	tests := []struct {
		name     string
		original int64
		rate     float64
		want     int64
	}{
		{"no discount", 10000, 0.0, 10000},
		{"10% off", 10000, 0.1, 9000},
		{"20% off", 10000, 0.2, 8000},
		{"50% off", 10000, 0.5, 5000},
		{"100% off", 10000, 1.0, 0},
		{
			"truncates fraction",
			999,
			0.1,
			899,
		}, // 999 * 900 / 1000 = 899.1 → 899
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := mustMoney(t, tt.original, domain.JPY)
			rate := mustRate(t, tt.rate)
			got, err := domain.CalculateDiscountedPrice(
				original,
				rate,
			)
			if err != nil {
				t.Fatalf("CalculateDiscountedPrice: %v", err)
			}
			want := mustMoney(t, tt.want, domain.JPY)
			if !got.Equals(want) {
				t.Errorf("want=%v got=%v", want, got)
			}
		})
	}
}

func TestNewDiscountRate_Invalid(t *testing.T) {
	tests := []struct {
		name string
		v    float64
	}{
		{"negative", -0.1},
		{"over 1", 1.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewDiscountRate(tt.v)
			if !errors.Is(err, domain.ErrInvalidDiscountRate) {
				t.Errorf("want ErrInvalidDiscountRate got %v", err)
			}
		})
	}
}

func TestMoney_AddSubtract_CurrencyMismatch(t *testing.T) {
	jpy := mustMoney(t, 1000, domain.JPY)
	usd := mustMoney(t, 10, domain.USD)

	if _, err := jpy.Add(usd); !errors.Is(
		err,
		domain.ErrCurrencyMismatch,
	) {
		t.Errorf(
			"Add cross-currency: want ErrCurrencyMismatch got %v",
			err,
		)
	}
	if _, err := jpy.Subtract(usd); !errors.Is(
		err,
		domain.ErrCurrencyMismatch,
	) {
		t.Errorf(
			"Subtract cross-currency: "+
				"want ErrCurrencyMismatch got %v",
			err,
		)
	}
}
