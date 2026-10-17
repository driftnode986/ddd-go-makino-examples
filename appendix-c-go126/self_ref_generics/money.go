// #region adder-currency
package money

import "fmt"

// Adder は「自分自身と加算できる」型の制約です。
// 制約の中で定義対象の型自身を参照することが
// Go 1.26 で許可されました。
type Adder[A Adder[A]] interface {
	Add(A) A
}

// Currency は通貨を区別するための phantom type です。
// 値を持たない struct を JPY / USD ごとに用意し、
// Money[JPY] と Money[USD] を別の型として扱わせます。
type Currency interface {
	currencyCode() string
}

type JPY struct{}

func (JPY) currencyCode() string { return "JPY" }

type USD struct{}

func (USD) currencyCode() string { return "USD" }

// #endregion

// #region money
// Money は通貨単位 C を型パラメータとして持つ値オブジェクトです。
type Money[C Currency] struct {
	amount int64
}

func NewMoney[C Currency](amount int64) Money[C] {
	return Money[C]{amount: amount}
}

// Add は同じ通貨単位の Money 同士でしか呼べません。
// Money[JPY].Add(Money[USD]) はコンパイルエラーになります。
func (m Money[C]) Add(other Money[C]) Money[C] {
	return Money[C]{amount: m.amount + other.amount}
}

func (m Money[C]) String() string {
	var c C
	return fmt.Sprintf("%d %s", m.amount, c.currencyCode())
}

// 自己参照制約に Money[JPY] が適合することを
// コンパイル時に確認します。
var _ Adder[Money[JPY]] = Money[JPY]{}

// #endregion
