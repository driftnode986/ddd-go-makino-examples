package postgres

import (
	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

func toRow(b *domain.Booking) BookingRow {
	return BookingRow{
		ID:           b.ID().String(),
		MemberID:     b.MemberID().String(),
		LessonID:     b.LessonID().String(),
		ScheduledAt:  b.ScheduledAt(),
		Status:       string(b.Status()),
		CancelReason: b.CancelReason(),
	}
}
