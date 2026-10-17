// #region detector
package domain

import (
	"context"
	"errors"
)

// ErrBookingConflict は予約時間の衝突を表すドメインエラー。
var ErrBookingConflict = errors.New("booking conflicts with existing reservation")

// ConflictDetector は同じトレーナーの予約時間衝突を検出するドメインサービス。
type ConflictDetector struct {
	bookings BookingFinder
}

// BookingFinder は時間帯検索の interface (リポジトリの一部を切り出した小さい interface)。
type BookingFinder interface {
	FindByTrainerAndPeriod(ctx context.Context, trainerID TrainerID, from, to Hour) ([]*Booking, error)
}

// NewConflictDetector はファクトリ関数。
func NewConflictDetector(bookings BookingFinder) *ConflictDetector {
	return &ConflictDetector{bookings: bookings}
}

// #endregion

// #region detect
// Detect は与えられた時間帯にトレーナーの予約衝突があれば ErrBookingConflict を返す。
func (d *ConflictDetector) Detect(
	ctx context.Context, trainerID TrainerID, from, to Hour,
) error {
	if !from.Before(to) {
		return ErrInvalidPeriod
	}
	hits, err := d.bookings.FindByTrainerAndPeriod(ctx, trainerID, from, to)
	if err != nil {
		return err
	}
	for _, b := range hits {
		if b.Status() == StatusCancelled {
			continue
		}
		if HoursOverlap(from, to, b.From(), b.To()) {
			return ErrBookingConflict
		}
	}
	return nil
}

// #endregion
