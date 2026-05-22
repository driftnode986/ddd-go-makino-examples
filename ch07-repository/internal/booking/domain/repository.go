package domain

import "context"

type Repository interface {
	Add(ctx context.Context, b *Booking) error
	Save(ctx context.Context, b *Booking) error
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
	Remove(ctx context.Context, id BookingID) error
}
