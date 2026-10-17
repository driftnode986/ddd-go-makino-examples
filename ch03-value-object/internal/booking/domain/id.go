// #region id-type
package domain

import (
	"errors"
	"fmt"
)

// ID は phantom type を使った型安全な ID。
// T には User, Booking, Lesson などのドメイン概念型を渡す。
type ID[T any] struct {
	value string
}

// ErrInvalidID は不正な ID 文字列を表すドメインエラー。
var ErrInvalidID = errors.New("invalid id")

// NewID は ID 文字列 (UUID 想定) を検証して ID[T] を返す。
func NewID[T any](s string) (ID[T], error) {
	if len(s) < 8 {
		return ID[T]{}, fmt.Errorf("%w: too short", ErrInvalidID)
	}
	return ID[T]{value: s}, nil
}

// String は内部値を返す (fmt.Stringer)。
func (id ID[T]) String() string { return id.value }

// Equals は同じ T 型同士の比較。
// 異なる T の ID は型エラーで弾かれる。
func (id ID[T]) Equals(
	other ID[T],
) bool {
	return id.value == other.value
}

// #endregion

// 以下、ドメイン用語を phantom type として宣言する。

// #region id-markers
// User は User エンティティを示すマーカー型。
type User struct{}

// Booking は Booking 集約を示すマーカー型。
type Booking struct{}

// Lesson は Lesson エンティティを示すマーカー型。
type Lesson struct{}

// UserID, BookingID, LessonID は型安全な ID の別名。
type (
	UserID    = ID[User]
	BookingID = ID[Booking]
	LessonID  = ID[Lesson]
)

// #endregion
