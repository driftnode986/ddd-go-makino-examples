package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/domain"
)

func TestNewID_RejectsTooShort(t *testing.T) {
	_, err := domain.NewID[domain.User]("short")
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("want ErrInvalidID, got %v", err)
	}
}

func TestUserID_Equals(t *testing.T) {
	a, _ := domain.NewID[domain.User]("user-1234")
	b, _ := domain.NewID[domain.User]("user-1234")
	c, _ := domain.NewID[domain.User]("user-9999")

	if !a.Equals(b) {
		t.Errorf("a and b should be equal")
	}
	if a.Equals(c) {
		t.Errorf("a and c should not be equal")
	}
}

// 以下は compile time の型安全を確認するためのコード例。
// 異なる phantom type の ID 同士を直接比較しようとすると compile error になる。
//
//	uid, _ := domain.NewID[domain.User]("user-1234")
//	bid, _ := domain.NewID[domain.Booking]("book-1234")
//	uid.Equals(bid)  // compile error: ID[User] と ID[Booking] は別の型
func TestPhantomType_PreventsTypeConfusion(t *testing.T) {
	uid, _ := domain.NewID[domain.User]("user-1234")
	bid, _ := domain.NewID[domain.Booking]("book-1234")

	// 同じ phantom type 同士は比較できる
	uid2, _ := domain.NewID[domain.User]("user-1234")
	if !uid.Equals(uid2) {
		t.Errorf("same phantom type IDs should be comparable")
	}

	// 内部値が偶然一致しても、phantom type が違えば別物として扱える
	// (Equals は同じ T のみ受け付けるため、bid を uid.Equals に渡そうとすると compile error)
	_ = bid
}
