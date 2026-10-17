package domain

import (
	"context"
	"iter"
)

type BookingRepository interface {
	Save(ctx context.Context, b *Booking) error
	FindByID(ctx context.Context, id BookingID) (*Booking, error)
	// 大量集約を遅延評価で返す
	AllByMember(ctx context.Context, memberID MemberID) iter.Seq2[*Booking, error]
}
