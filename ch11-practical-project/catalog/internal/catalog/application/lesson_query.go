// catalog/internal/catalog/application/lesson_query.go
package application

import (
	"context"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

type LessonQuery struct {
	repo domain.LessonRepository
}

func NewLessonQuery(repo domain.LessonRepository) *LessonQuery {
	return &LessonQuery{repo: repo}
}

// FindLessonView は Catalog BC が他 BC に公開する公表された言語 (LessonView) を返します。
// Booking BC はこのメソッドのシグネチャを application.LessonQueryPort interface で受けます。
func (q *LessonQuery) FindLessonView(ctx context.Context,
	idStr string) (dto.LessonView, error) {
	id, err := domain.LessonIDFrom(idStr)
	if err != nil {
		return dto.LessonView{}, fmt.Errorf("invalid id: %w", err)
	}
	l, err := q.repo.FindByID(ctx, id)
	if err != nil {
		return dto.LessonView{}, err
	}
	return dto.LessonView{
		ID:        l.ID().String(),
		Name:      l.Name(),
		Capacity:  l.Capacity(),
		PriceJPY:  l.PriceJPY(),
		StartsAt:  l.StartsAt(),
		TrainerID: l.TrainerID(),
	}, nil
}
