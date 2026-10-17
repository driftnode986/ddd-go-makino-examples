package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type LessonID struct{ value string }

func NewLessonID() LessonID {
	return LessonID{value: uuid.NewString()}
}

func LessonIDFrom(s string) (LessonID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return LessonID{}, errors.New("invalid lesson id")
	}
	return LessonID{value: s}, nil
}

func (id LessonID) String() string { return id.value }

type Lesson struct {
	id        LessonID
	name      string
	capacity  int
	priceJPY  int
	startsAt  time.Time
	trainerID string
}

func NewLesson(name string, capacity, priceJPY int,
	startsAt time.Time, trainerID string) (*Lesson, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	if priceJPY < 0 {
		return nil, errors.New("price must be zero or positive")
	}
	return &Lesson{
		id:        NewLessonID(),
		name:      name,
		capacity:  capacity,
		priceJPY:  priceJPY,
		startsAt:  startsAt,
		trainerID: trainerID,
	}, nil
}

func (l *Lesson) ID() LessonID        { return l.id }
func (l *Lesson) Name() string        { return l.name }
func (l *Lesson) Capacity() int       { return l.capacity }
func (l *Lesson) PriceJPY() int       { return l.priceJPY }
func (l *Lesson) StartsAt() time.Time { return l.startsAt }
func (l *Lesson) TrainerID() string   { return l.trainerID }
