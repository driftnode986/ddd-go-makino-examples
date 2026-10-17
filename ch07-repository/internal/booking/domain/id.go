package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidID = errors.New("invalid id")

type BookingID struct {
	value string
}

func NewBookingID(s string) (BookingID, error) {
	if s == "" {
		return BookingID{}, fmt.Errorf(
			"booking id: %w",
			ErrInvalidID,
		)
	}
	return BookingID{value: s}, nil
}

func ParseBookingID(
	s string,
) (BookingID, error) {
	return NewBookingID(s)
}

func (id BookingID) String() string { return id.value }

type MemberID struct {
	value string
}

func NewMemberID(s string) (MemberID, error) {
	if s == "" {
		return MemberID{}, fmt.Errorf("member id: %w", ErrInvalidID)
	}
	return MemberID{value: s}, nil
}

func ParseMemberID(
	s string,
) (MemberID, error) {
	return NewMemberID(s)
}

func (id MemberID) String() string { return id.value }

type LessonID struct {
	value string
}

func NewLessonID(s string) (LessonID, error) {
	if s == "" {
		return LessonID{}, fmt.Errorf("lesson id: %w", ErrInvalidID)
	}
	return LessonID{value: s}, nil
}

func ParseLessonID(
	s string,
) (LessonID, error) {
	return NewLessonID(s)
}

func (id LessonID) String() string { return id.value }
