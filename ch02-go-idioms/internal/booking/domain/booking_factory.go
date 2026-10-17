package domain

import (
	"errors"
	"time"

	shared "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/shared/domain"
)

var (
	ErrInvalidBooking = errors.New(
		"member, lesson and scheduled time are required",
	)
	ErrBookingCannotReschedule = errors.New(
		"booking cannot be rescheduled in current state",
	)
	ErrBookingCannotComplete = errors.New(
		"booking cannot be completed in current state",
	)
)

// NewBooking はファクトリ関数。新しい BookingID を振り、確定前 (pending) の予約を作る。
func NewBooking(
	memberID MemberID,
	lessonID LessonID,
	scheduledAt time.Time,
) (*Booking, error) {
	if memberID == "" || lessonID == "" || scheduledAt.IsZero() {
		return nil, ErrInvalidBooking
	}
	return &Booking{
		id:          shared.NewBookingID(),
		memberID:    memberID,
		lessonID:    lessonID,
		scheduledAt: scheduledAt,
		status:      StatusPending,
	}, nil
}
