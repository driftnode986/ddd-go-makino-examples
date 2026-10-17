package infrastructure_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/infrastructure"
)

func TestMemoryUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := infrastructure.NewMemoryUserRepository()
	email, _ := domain.NewEmail("taro@example.com")
	u, _ := domain.NewUser("usr_1", email, "Taro", time.Now())

	if err := repo.Add(ctx, u); err != nil {
		t.Fatal(err)
	}
	if err := repo.Add(ctx, u); !errors.Is(err, infrastructure.ErrDuplicateUser) {
		t.Errorf("second Add err = %v", err)
	}
	if err := repo.Save(ctx, u); err != nil {
		t.Errorf("Save: %v", err)
	}
	if got, err := repo.FindByID(ctx, "usr_1"); err != nil || got != u {
		t.Errorf("FindByID = %v, %v", got, err)
	}
	other, _ := domain.NewEmail("x@example.com")
	if _, err := repo.FindByEmail(ctx, other); !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("FindByEmail(other) err = %v", err)
	}
}

func TestRandomIDGenerator(t *testing.T) {
	g := infrastructure.RandomIDGenerator{}
	a, b := g.New(), g.New()
	if a == b || !strings.HasPrefix(a.String(), "usr_") || len(a) != len("usr_")+16 {
		t.Errorf("ids = %q, %q", a, b)
	}
}
