package money_test

import (
	"testing"

	money "github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/self_ref_generics"
)

func TestMoneyAdd_SameCurrency(t *testing.T) {
	a := money.NewMoney[money.JPY](1000)
	b := money.NewMoney[money.JPY](500)
	if got := a.Add(b).String(); got != "1500 JPY" {
		t.Errorf("got %q, want %q", got, "1500 JPY")
	}
}

// sum は Adder 制約を満たす型なら何でも足し合わせる。
func sum[A money.Adder[A]](first A, rest ...A) A {
	total := first
	for _, v := range rest {
		total = total.Add(v)
	}
	return total
}

func TestAdderConstraint(t *testing.T) {
	got := sum(money.NewMoney[money.USD](1), money.NewMoney[money.USD](2), money.NewMoney[money.USD](3))
	if got.String() != "6 USD" {
		t.Errorf("got %q, want %q", got.String(), "6 USD")
	}
}

// 次のコードはコンパイルエラーになることを意図的に示しています。
// money.NewMoney[money.JPY](100).Add(money.NewMoney[money.USD](100))
// → cannot use ... (Money[USD]) as Money[JPY] value in argument
