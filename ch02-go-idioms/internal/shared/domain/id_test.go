package domain_test

import (
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/shared/domain"
	"github.com/google/uuid"
)

func TestNewID_GeneratesDistinctUUIDs(t *testing.T) {
	a := domain.NewUserID()
	b := domain.NewUserID()
	if a == b {
		t.Fatalf("two generated IDs should differ: %s", a)
	}
	if _, err := uuid.Parse(a.String()); err != nil {
		t.Fatalf(
			"String() should be a UUID: %q (%v)",
			a.String(),
			err,
		)
	}
}

func TestNewBookingID_IsUUID(t *testing.T) {
	id := domain.NewBookingID()
	if _, err := uuid.Parse(id.String()); err != nil {
		t.Fatalf(
			"String() should be a UUID: %q (%v)",
			id.String(),
			err,
		)
	}
}

// NewID[T] は任意のタグ型で使える。UserID と BookingID は別の型なので、
// var u domain.UserID = domain.NewBookingID() はコンパイルエラーになる。
func TestNewID_WithOwnTag(t *testing.T) {
	type lessonTag struct{}
	id := domain.NewID[lessonTag]()
	if id.String() == domain.NewID[lessonTag]().String() {
		t.Fatal("two generated IDs should differ")
	}
}
