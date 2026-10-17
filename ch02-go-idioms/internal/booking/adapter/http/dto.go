package http

// CreateBookingRequest は HTTP リクエストの DTO。公開フィールドで OK
type CreateBookingRequest struct {
	MemberID    string `json:"member_id"`
	LessonID    string `json:"lesson_id"`
	ScheduledAt string `json:"scheduled_at"`
}
