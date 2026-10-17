package domain

import "time"

type Booking struct {
	id          BookingID
	memberID    MemberID
	lessonID    LessonID
	scheduledAt time.Time
	status      BookingStatus
}

func (b *Booking) Confirm() error { /* ... */ return nil }

func (b *Booking) Cancel(
	reason string,
) error { /* ... */
	return nil
}

func (b *Booking) Reschedule(
	newAt time.Time,
) error { /* ... */
	return nil
}

func (b *Booking) MarkCompleted() error { /* ... */ return nil }
func (b *Booking) ID() BookingID        { return b.id }

func (b *Booking) Status() BookingStatus { return b.status }
