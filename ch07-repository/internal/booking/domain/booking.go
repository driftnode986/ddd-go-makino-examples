package domain

import "time"

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	id          BookingID
	memberID    MemberID
	lessonID    LessonID
	status      BookingStatus
	scheduledAt time.Time
	version     int
}

func NewBooking(id BookingID, memberID MemberID, lessonID LessonID, scheduledAt time.Time) *Booking {
	return &Booking{
		id:          id,
		memberID:    memberID,
		lessonID:    lessonID,
		status:      StatusPending,
		scheduledAt: scheduledAt,
		version:     1,
	}
}

func ReconstructBooking(
	id BookingID,
	memberID MemberID,
	lessonID LessonID,
	status BookingStatus,
	scheduledAt time.Time,
	version int,
) *Booking {
	return &Booking{
		id:          id,
		memberID:    memberID,
		lessonID:    lessonID,
		status:      status,
		scheduledAt: scheduledAt,
		version:     version,
	}
}

func (b *Booking) ID() BookingID          { return b.id }
func (b *Booking) MemberID() MemberID     { return b.memberID }
func (b *Booking) LessonID() LessonID     { return b.lessonID }
func (b *Booking) Status() BookingStatus  { return b.status }
func (b *Booking) ScheduledAt() time.Time { return b.scheduledAt }
func (b *Booking) Version() int           { return b.version }

func (b *Booking) Confirm() error {
	if b.status != StatusPending {
		return ErrCannotCancel
	}
	b.status = StatusConfirmed
	b.version++
	return nil
}

func (b *Booking) Cancel() error {
	if b.status == StatusCancelled {
		return ErrCannotCancel
	}
	b.status = StatusCancelled
	b.version++
	return nil
}
