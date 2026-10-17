package memory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/infrastructure/memory"
)

func TestMemberRepository_SaveAndFind(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemberRepository()
	email, _ := domain.NewEmail("taro@example.com")
	m, err := domain.NewMember(email, "Taro")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, m); err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindByID(ctx, m.ID())
	if err != nil || got.Name() != "Taro" {
		t.Fatalf("FindByID = %v, %v", got, err)
	}
	if _, err := repo.FindByID(ctx, domain.NewMemberID()); !errors.Is(err, domain.ErrMemberNotFound) {
		t.Errorf("err = %v, want ErrMemberNotFound", err)
	}
}
