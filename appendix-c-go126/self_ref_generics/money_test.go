package money_test

import (
	"testing"

	money "github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/self_ref_generics"
)

func TestMoneyAdd_SameCurrency(t *testing.T) {
	a := money.NewMoney[money.JPY](1000)
	b := money.NewMoney[money.JPY](500)
	got := a.Add(b)
	if got.Amount() != 1500 {
		t.Errorf("got %d, want 1500", got.Amount())
	}
}

// 次のコードはコンパイルエラーになることを意図的に示しています。
// money.NewMoney[money.JPY](100).Add(money.NewMoney[money.USD](100))
// → cannot use ... (Money[USD]) as Money[JPY] value in argument
