package domain

import (
	"errors"

	"github.com/google/uuid"
)

type BookingID struct{ value string }

func NewBookingID() BookingID {
	return BookingID{value: uuid.NewString()}
}

func BookingIDFrom(s string) (BookingID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return BookingID{}, errors.New("invalid booking id")
	}
	return BookingID{value: s}, nil
}

func (id BookingID) String() string { return id.value }

type MemberID struct{ value string }

func NewMemberID() MemberID { return MemberID{value: uuid.NewString()} }

func MemberIDFrom(s string) (MemberID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return MemberID{}, errors.New("invalid member id")
	}
	return MemberID{value: s}, nil
}

func (id MemberID) String() string { return id.value }

type LessonID struct{ value string }

func NewLessonID() LessonID { return LessonID{value: uuid.NewString()} }

func LessonIDFrom(s string) (LessonID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return LessonID{}, errors.New("invalid lesson id")
	}
	return LessonID{value: s}, nil
}

func (id LessonID) String() string { return id.value }
