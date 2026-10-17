package application

import (
	"context"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// ProcessMemberBookings は会員の予約を 1 件ずつ取り出して process に渡す。
func ProcessMemberBookings(
	ctx context.Context,
	repo domain.BookingRepository,
	memberID domain.MemberID,
	process func(*domain.Booking) error,
) error {
	// #region range-bookings
	// 集約を 1 件ずつストリーミング処理する
	for b, err := range repo.AllByMember(ctx, memberID) {
		if err != nil {
			return fmt.Errorf("iterate bookings: %w", err)
		}
		if err := process(b); err != nil {
			return err
		}
	}
	// #endregion
	return nil
}
