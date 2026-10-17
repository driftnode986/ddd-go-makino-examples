package domain

import (
	"context"
	"errors"
)

var ErrBookingNotFound = errors.New("booking not found")

type BookingRepository interface {
	Save(ctx context.Context, b *Booking) error
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
	ListByMember(
		ctx context.Context,
		memberID MemberID,
	) ([]*Booking, error)
}
