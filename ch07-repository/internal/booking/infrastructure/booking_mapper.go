package infrastructure

import (
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
)

func toBookingRow(b *domain.Booking) bookingRow {
	return bookingRow{
		ID:          b.ID().String(),
		MemberID:    b.MemberID().String(),
		LessonID:    b.LessonID().String(),
		Status:      string(b.Status()),
		ScheduledAt: b.ScheduledAt(),
		Version:     b.Version(),
	}
}

func toBookingDomain(r bookingRow) (*domain.Booking, error) {
	id, err := domain.ParseBookingID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid booking id %q: %w", r.ID, err)
	}
	memberID, err := domain.ParseMemberID(r.MemberID)
	if err != nil {
		return nil, fmt.Errorf("invalid member id %q: %w", r.MemberID, err)
	}
	lessonID, err := domain.ParseLessonID(r.LessonID)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson id %q: %w", r.LessonID, err)
	}
	status := domain.BookingStatus(r.Status)
	return domain.ReconstructBooking(
		id, memberID, lessonID, status, r.ScheduledAt, r.Version,
	), nil
}
