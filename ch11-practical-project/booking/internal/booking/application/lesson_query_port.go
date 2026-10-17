package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

// LessonQueryPort は Catalog BC への参照ポートです。
// Booking BC は catalog/internal/... を import せず、shared/dto を介して連携します。
type LessonQueryPort interface {
	FindLessonView(ctx context.Context, lessonID string) (dto.LessonView, error)
}
