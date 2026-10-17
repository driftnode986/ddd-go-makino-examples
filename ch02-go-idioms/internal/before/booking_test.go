package before_test

import (
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/before"
)

func TestBooking_CancelBySetters(t *testing.T) {
	b := before.NewBooking("b-1", "m-1")
	// キャンセルするには、呼び出し側が 2 つの Setter を正しい順で呼ぶ必要がある
	b.SetStatus("cancelled")
	b.SetCancelReason("schedule changed")
	if b.Status != "cancelled" ||
		b.GetCancelReason() != "schedule changed" {
		t.Fatalf("got %+v", b)
	}
}

func TestBooking_NoRuleStopsInvalidTransition(t *testing.T) {
	b := before.NewBooking("b-1", "m-1")
	b.SetStatus("completed")
	b.SetStatus("pending") // 完了後に戻すことも止められない
	b.Status = "whatever"  // 公開フィールドなので任意の文字列も入る
	if b.Status != "whatever" {
		t.Fatalf("got %q", b.Status)
	}
}

func TestNewBooking_PanicsOnEmptyID(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("want panic for empty id")
		}
	}()
	before.NewBooking("", "m-1")
}
