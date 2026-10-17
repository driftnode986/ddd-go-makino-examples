package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/domain"
)

func TestNewMoney_RejectsNegativeAmount(t *testing.T) {
	_, err := domain.NewMoney(-1, domain.JPY)
	if !errors.Is(err, domain.ErrNegativeAmount) {
		t.Fatalf("want ErrNegativeAmount, got %v", err)
	}
}

func TestMoney_Add_RequiresSameCurrency(t *testing.T) {
	yen, _ := domain.NewMoney(1000, domain.JPY)
	usd, _ := domain.NewMoney(10, domain.USD)

	_, err := yen.Add(usd)
	if !errors.Is(err, domain.ErrCurrencyMismatch) {
		t.Fatalf("want ErrCurrencyMismatch, got %v", err)
	}
}

// #region immutability
func TestMoney_Add_PreservesImmutability(t *testing.T) {
	a, _ := domain.NewMoney(1000, domain.JPY)
	b, _ := domain.NewMoney(500, domain.JPY)

	sum, err := a.Add(b)
	if err != nil {
		t.Fatal(err)
	}
	if sum.Amount() != 1500 {
		t.Errorf("sum: want 1500, got %d", sum.Amount())
	}
	// a と b は変化していない (不変性)
	if a.Amount() != 1000 {
		t.Errorf("a was mutated: want 1000, got %d", a.Amount())
	}
	if b.Amount() != 500 {
		t.Errorf("b was mutated: want 500, got %d", b.Amount())
	}
}

// #endregion

func TestMoney_Equals(t *testing.T) {
	a, _ := domain.NewMoney(1000, domain.JPY)
	b, _ := domain.NewMoney(1000, domain.JPY)
	c, _ := domain.NewMoney(1000, domain.USD)

	if !a.Equals(b) {
		t.Errorf("a and b should be equal")
	}
	if a.Equals(c) {
		t.Errorf("a (JPY) and c (USD) should not be equal")
	}
}
