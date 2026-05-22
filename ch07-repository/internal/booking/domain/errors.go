package domain

import "errors"

var (
	ErrNotFound        = errors.New("booking not found")
	ErrDuplicate       = errors.New("booking already exists")
	ErrVersionConflict = errors.New("booking version conflict")
	ErrCannotCancel    = errors.New("booking cannot be cancelled in current state")
)
