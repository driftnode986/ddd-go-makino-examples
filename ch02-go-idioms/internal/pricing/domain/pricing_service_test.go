package pricing_test

import (
	"errors"
	"testing"
	"time"

	pricing "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/pricing/domain"
)

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func yen(t *testing.T, amount int64) pricing.Money {
	t.Helper()
	m, err := pricing.NewMoney(amount, "JPY")
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestPricingService_CalculateFinalPrice(t *testing.T) {
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	svc := pricing.NewPricingService(fixedClock{t: now})

	t.Run("valid coupon", func(t *testing.T) {
		coupon := pricing.NewCoupon(yen(t, 500), now.Add(time.Hour))
		got, err := svc.CalculateFinalPrice(yen(t, 3000), coupon)
		if err != nil {
			t.Fatal(err)
		}
		if got.Amount() != 2500 {
			t.Fatalf("want 2500, got %d", got.Amount())
		}
	})
	t.Run("expired coupon", func(t *testing.T) {
		coupon := pricing.NewCoupon(yen(t, 500), now.Add(-time.Hour))
		if _, err := svc.CalculateFinalPrice(yen(t, 3000), coupon); err == nil {
			t.Fatal("want error for expired coupon")
		}
	})
	t.Run("discount larger than price", func(t *testing.T) {
		coupon := pricing.NewCoupon(yen(t, 5000), now.Add(time.Hour))
		_, err := svc.CalculateFinalPrice(yen(t, 3000), coupon)
		if !errors.Is(err, pricing.ErrNegativeAmount) {
			t.Fatalf("want ErrNegativeAmount, got %v", err)
		}
	})
}

func TestMoney_Subtract_CurrencyMismatch(t *testing.T) {
	usd, _ := pricing.NewMoney(10, "USD")
	if _, err := yen(t, 100).Subtract(usd); !errors.Is(err, pricing.ErrCurrencyMismatch) {
		t.Fatalf("want ErrCurrencyMismatch, got %v", err)
	}
}
