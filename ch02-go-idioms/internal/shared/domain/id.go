package domain

import "github.com/google/uuid"

// ID[T] は型パラメータ T で識別される ID 型
// UserID と BookingID を取り違える型エラーを compile time で防ぐ
type ID[T any] struct {
	value uuid.UUID
}

func NewID[T any]() ID[T] {
	return ID[T]{value: uuid.New()}
}

func (id ID[T]) String() string { return id.value.String() }

// 型エイリアスで業務用語に対応
type (
	userTag    struct{}
	bookingTag struct{}
)

type UserID = ID[userTag]
type BookingID = ID[bookingTag]
