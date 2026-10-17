package domain

import (
	"context"
	"errors"
)

var ErrLessonNotFound = errors.New("lesson not found")

type LessonRepository interface {
	Save(ctx context.Context, l *Lesson) error
	FindByID(ctx context.Context, id LessonID) (*Lesson, error)
	List(ctx context.Context) ([]*Lesson, error)
}
