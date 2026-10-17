package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/infrastructure/memory"
)

func TestLessonRepository_SaveFindList(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewLessonRepository()
	l, err := domain.NewLesson(
		"Yoga",
		10,
		2000,
		time.Now(),
		"trainer-1",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, l); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindByID(ctx, l.ID())
	if err != nil || got.Name() != "Yoga" {
		t.Fatalf("FindByID = %v, %v", got, err)
	}
	all, err := repo.List(ctx)
	if err != nil || len(all) != 1 {
		t.Errorf("List = %d, %v", len(all), err)
	}
	if _, err := repo.FindByID(
		ctx,
		domain.NewLessonID(),
	); !errors.Is(
		err,
		domain.ErrLessonNotFound,
	) {
		t.Errorf("err = %v, want ErrLessonNotFound", err)
	}
}
