package domain

import (
	"fmt"
	"time"
)

type ScheduleConflictError struct {
	ExistingBookingID BookingID
	RequestedAt       time.Time
}

func (e *ScheduleConflictError) Error() string {
	return fmt.Sprintf("schedule conflict with %s at %v",
		e.ExistingBookingID, e.RequestedAt)
}

// Is メソッドで errors.Is(err, ErrScheduleConflict) を有効にする
func (e *ScheduleConflictError) Is(target error) bool {
	return target == ErrScheduleConflict
}
