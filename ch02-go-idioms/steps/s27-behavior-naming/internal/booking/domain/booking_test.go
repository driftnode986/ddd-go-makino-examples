package domain

import (
	"testing"

	canonical "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// Setter は検証なしで状態を書き換える。キャンセル理由の記録は呼び出し側の責任になる。
func TestSetter_ChangesStateWithoutRules(t *testing.T) {
	b := &Booking{status: canonical.StatusCompleted}
	b.SetStatus(canonical.StatusCancelled) // 完了済みでも止められない
	if b.status != canonical.StatusCancelled ||
		b.cancelReason != "" {
		t.Fatalf("got status=%s reason=%q", b.status, b.cancelReason)
	}
	b.SetCancelReason("late")
	if b.cancelReason != "late" {
		t.Fatalf("reason: got %q", b.cancelReason)
	}
}

// ビヘイビア側は中身を省略した版なので、状態を変えずに nil を返す。
func TestBehaviorMethods_AreStubs(t *testing.T) {
	b := &Booking{status: canonical.StatusPending}
	if b.Confirm() != nil || b.Cancel("x") != nil ||
		b.MarkCompleted() != nil {
		t.Fatal("stubs should return nil")
	}
	if b.status != canonical.StatusPending {
		t.Fatalf("status: got %s", b.status)
	}
}
