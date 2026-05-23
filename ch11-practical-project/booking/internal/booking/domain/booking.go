// internal/booking/domain/booking.go
package domain

import (
	"errors"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	id        BookingID
	memberID  MemberID
	lessonID  LessonID
	status    BookingStatus
	version   int
	createdAt time.Time
}

func NewBooking(memberID MemberID, lessonID LessonID, now time.Time) *Booking {
	return &Booking{
		id:        NewBookingID(),
		memberID:  memberID,
		lessonID:  lessonID,
		status:    BookingStatusPending,
		version:   1,
		createdAt: now,
	}
}

func RestoreBooking(id BookingID, memberID MemberID, lessonID LessonID,
	status BookingStatus, version int, createdAt time.Time) *Booking {
	return &Booking{id: id, memberID: memberID, lessonID: lessonID,
		status: status, version: version, createdAt: createdAt}
}

func (b *Booking) ID() BookingID         { return b.id }
func (b *Booking) MemberID() MemberID    { return b.memberID }
func (b *Booking) LessonID() LessonID    { return b.lessonID }
func (b *Booking) Status() BookingStatus { return b.status }
func (b *Booking) Version() int          { return b.version }
func (b *Booking) CreatedAt() time.Time  { return b.createdAt }

func (b *Booking) Confirm() error {
	if b.status != BookingStatusPending {
		return errors.New("cannot confirm a non-pending booking")
	}
	b.status = BookingStatusConfirmed
	b.version++
	return nil
}

func (b *Booking) Cancel() error {
	if b.status == BookingStatusCancelled {
		return errors.New("booking is already cancelled")
	}
	b.status = BookingStatusCancelled
	b.version++
	return nil
}
