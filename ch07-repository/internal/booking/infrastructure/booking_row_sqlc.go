package infrastructure

import "github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure/sqlcgen"

// rowToBookingRow は sqlc が生成した行の型 (sqlcgen.Booking) を
// Persistence Model (bookingRow) に詰め替える。ListByMember が使う。
func rowToBookingRow(row sqlcgen.Booking) bookingRow {
	return bookingRow{
		ID:          row.ID,
		MemberID:    row.MemberID,
		LessonID:    row.LessonID,
		Status:      row.Status,
		ScheduledAt: row.ScheduledAt.Time,
		Version:     int(row.Version),
	}
}
