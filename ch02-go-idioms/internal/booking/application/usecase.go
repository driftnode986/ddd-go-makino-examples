package application

import (
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// BookLessonInput はレッスン予約のユースケースへの入力。
type BookLessonInput struct {
	MemberID    domain.MemberID
	LessonID    domain.LessonID
	ScheduledAt time.Time
}

// BookLessonUseCase はレッスンを予約するアプリケーションサービス (第8章で詳しく扱う)。
type BookLessonUseCase struct {
	repo domain.BookingRepository
}

// NewBookLessonUseCase はリポジトリを受け取ってユースケースを作る。
func NewBookLessonUseCase(repo domain.BookingRepository) *BookLessonUseCase {
	return &BookLessonUseCase{repo: repo}
}
