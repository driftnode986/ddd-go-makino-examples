package domain

import "context"

// BookingRepository は集約の永続化を抽象化する。
// 小さい interface に保つ
type BookingRepository interface {
	Save(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
	Remove(ctx context.Context, id BookingID) error
}
