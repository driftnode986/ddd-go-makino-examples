// #region id-type
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

// #endregion

// #region id-methods
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

// #endregion

// 以下、ドメイン用語を phantom type のマーカー型として宣言する。
// マーカー型名には bookingTag のように小文字始まりまたは Tag サフィックスを使い、
// エンティティ本体 (Booking, Member, Lesson) と名前衝突を避ける。

// #region id-markers
// BookingTag は Booking エンティティ用のマーカー型。
type BookingTag struct{}

// MemberTag は Member エンティティ用のマーカー型。
type MemberTag struct{}

// LessonTag は Lesson エンティティ用のマーカー型。
type LessonTag struct{}

// BookingID, MemberID, LessonID は型安全な ID の別名。
type (
	BookingID = ID[BookingTag]
	MemberID  = ID[MemberTag]
	LessonID  = ID[LessonTag]
)

// #endregion
