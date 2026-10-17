// Package application は 2.3 節の「errors.Is / errors.As で判別する」版のユースケース。
// 2.12 節で log/slog を使う形に書き換えた完成形は
// ch02-go-idioms/internal/booking/application にある。
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

// BookLessonUseCase はレッスンを予約するアプリケーションサービス。
type BookLessonUseCase struct {
	repo domain.BookingRepository
}

// NewBookLessonUseCase はリポジトリを受け取ってユースケースを作る。
func NewBookLessonUseCase(
	repo domain.BookingRepository,
) *BookLessonUseCase {
	return &BookLessonUseCase{repo: repo}
}
