package domain

import "errors"

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingAlreadyConfirmed = errors.New(
		"booking already confirmed",
	)
	ErrBookingCannotCancel = errors.New(
		"booking cannot be cancelled in current state",
	)
	ErrScheduleConflict = errors.New("schedule conflict")
)
