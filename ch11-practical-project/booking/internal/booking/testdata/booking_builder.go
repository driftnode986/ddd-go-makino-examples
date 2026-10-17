// #region builder
package testdata

import (
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
)

type BookingBuilder struct {
	id        domain.BookingID
	memberID  domain.MemberID
	lessonID  domain.LessonID
	status    domain.BookingStatus
	version   int
	createdAt time.Time
}

func NewBookingBuilder() *BookingBuilder {
	return &BookingBuilder{
		id:        domain.NewBookingID(),
		memberID:  domain.NewMemberID(),
		lessonID:  domain.NewLessonID(),
		status:    domain.BookingStatusPending,
		version:   1,
		createdAt: time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC),
	}
}

// #endregion

// #region builder-methods
func (b *BookingBuilder) WithStatus(
	s domain.BookingStatus,
) *BookingBuilder {
	b.status = s
	return b
}

func (b *BookingBuilder) WithMember(
	id domain.MemberID,
) *BookingBuilder {
	b.memberID = id
	return b
}

func (b *BookingBuilder) Build() *domain.Booking {
	return domain.RestoreBooking(b.id, b.memberID, b.lessonID,
		b.status, b.version, b.createdAt)
}

// #endregion
