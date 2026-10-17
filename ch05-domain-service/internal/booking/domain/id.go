package domain

import (
	"errors"
	"fmt"
)

// ID は phantom type を使った型安全な ID。
type ID[T any] struct {
	value string
}

// ErrInvalidID は不正な ID 文字列を表すドメインエラー。
var ErrInvalidID = errors.New("invalid id")

// NewID は ID 文字列を検証して ID[T] を返す。
func NewID[T any](s string) (ID[T], error) {
	if len(s) < 8 {
		return ID[T]{}, fmt.Errorf("%w: too short", ErrInvalidID)
	}
	return ID[T]{value: s}, nil
}

// String は内部値を返す。
func (id ID[T]) String() string { return id.value }

// Equals は同じ T 型同士の比較。
func (id ID[T]) Equals(
	other ID[T],
) bool {
	return id.value == other.value
}

// IsZero はゼロ値かどうかを判定する。
func (id ID[T]) IsZero() bool { return id.value == "" }

// マーカー型。

// BookingTag は Booking エンティティ用のマーカー型。
type BookingTag struct{}

// MemberTag は Member エンティティ用のマーカー型。
type MemberTag struct{}

// LessonTag は Lesson エンティティ用のマーカー型。
type LessonTag struct{}

// TrainerTag は Trainer エンティティ用のマーカー型。
type TrainerTag struct{}

// BookingID, MemberID, LessonID, TrainerID は型安全な ID の別名。
type (
	BookingID = ID[BookingTag]
	MemberID  = ID[MemberTag]
	LessonID  = ID[LessonTag]
	TrainerID = ID[TrainerTag]
)
