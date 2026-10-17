package postgres

import "time"

// BookingRow は DB スキーマに対応する Persistence Model
type BookingRow struct {
	ID          string
	MemberID    string
	LessonID    string
	ScheduledAt time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
