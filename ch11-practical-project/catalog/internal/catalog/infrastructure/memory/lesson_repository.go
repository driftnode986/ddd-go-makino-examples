// catalog/internal/catalog/infrastructure/memory/lesson_repository.go
package memory

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
)

type LessonRepository struct {
	mu      sync.RWMutex
	lessons map[domain.LessonID]*domain.Lesson
}

func NewLessonRepository() *LessonRepository {
	return &LessonRepository{lessons: map[domain.LessonID]*domain.Lesson{}}
}

func (r *LessonRepository) Save(_ context.Context, l *domain.Lesson) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lessons[l.ID()] = l
	return nil
}

func (r *LessonRepository) FindByID(_ context.Context,
	id domain.LessonID) (*domain.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.lessons[id]
	if !ok {
		return nil, domain.ErrLessonNotFound
	}
	return l, nil
}

func (r *LessonRepository) List(_ context.Context) ([]*domain.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Lesson, 0, len(r.lessons))
	for _, l := range r.lessons {
		out = append(out, l)
	}
	return out, nil
}
