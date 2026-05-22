package infrastructure

import "time"

// bookingRow は bookings テーブルの 1 行を表す Persistence Model。
// ドメインモデルとは別構造体にし、ORM/SQL の詳細をここに閉じ込める。
type bookingRow struct {
	ID          string
	MemberID    string
	LessonID    string
	Status      string
	ScheduledAt time.Time
	Version     int
}
