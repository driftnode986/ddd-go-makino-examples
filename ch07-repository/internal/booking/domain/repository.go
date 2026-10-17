package domain

import "context"

// Repository は Booking 集約のリポジトリ。
// ドメイン層に interface を置き、実装は infrastructure 層に書く。
type Repository interface {
	Add(ctx context.Context, b *Booking) error
	Save(ctx context.Context, b *Booking) error
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
	Remove(ctx context.Context, id BookingID) error
}
