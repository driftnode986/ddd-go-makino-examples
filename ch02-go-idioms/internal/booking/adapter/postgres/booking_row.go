package postgres

import "time"

type BookingRow struct {
	ID           string    `db:"id"`
	MemberID     string    `db:"member_id"`
	LessonID     string    `db:"lesson_id"`
	ScheduledAt  time.Time `db:"scheduled_at"`
	Status       string    `db:"status"`
	CancelReason string    `db:"cancel_reason"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
