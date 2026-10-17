package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/infrastructure/memory"
)

func TestLessonQuery_FindLessonView(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewLessonRepository()
	startsAt := time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC)
	l, err := domain.NewLesson(
		"Yoga",
		12,
		2500,
		startsAt,
		"trainer-1",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, l); err != nil {
		t.Fatal(err)
	}

	view, err := application.NewLessonQuery(repo).
		FindLessonView(ctx, l.ID().String())
	if err != nil {
		t.Fatalf("FindLessonView: %v", err)
	}
	if view.ID != l.ID().String() || view.Name != "Yoga" ||
		view.Capacity != 12 || view.PriceJPY != 2500 ||
		!view.StartsAt.Equal(startsAt) ||
		view.TrainerID != "trainer-1" {
		t.Errorf("view = %+v", view)
	}
	if _, err := application.NewLessonQuery(repo).FindLessonView(
		ctx,
		"not-uuid",
	); err == nil {
		t.Error("FindLessonView(not-uuid): want error, got nil")
	}
}
