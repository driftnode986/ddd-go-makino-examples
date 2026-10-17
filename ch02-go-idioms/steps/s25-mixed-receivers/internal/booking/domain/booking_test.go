package domain

import "testing"

// 混在させても、アドレスを取れる変数ならポインタレシーバのメソッドも呼べる。
// 呼べなくなるのは、マップの要素や関数の戻り値のようにアドレスを取れない値。
func TestMixedReceivers_AddressableValueCanCallPointerMethod(
	t *testing.T,
) {
	var b Booking
	if err := b.Confirm(); err != nil { // (&b).Confirm() として呼ばれる
		t.Fatal(err)
	}
	var zero BookingID
	if b.ID() != zero {
		t.Fatal("ID: want the zero value")
	}
}

// 値の Booking はポインタレシーバのメソッドを持たないので、Confirm を含む
// interface を満たさない (var _ confirmer = Booking{} はコンパイルエラー)。
type confirmer interface{ Confirm() error }

func TestMixedReceivers_OnlyPointerSatisfiesInterface(t *testing.T) {
	var c confirmer = &Booking{}
	if err := c.Confirm(); err != nil {
		t.Fatal(err)
	}
}
