package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

// fakeTx は fn をそのまま呼び、呼ばれた回数を数える。
type fakeTx struct{ calls int }

func (f *fakeTx) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.calls++
	return fn(ctx)
}

type fixedID domain.UserID

func (id fixedID) New() domain.UserID { return domain.UserID(id) }

func TestRegisterUserUseCase_AddsUserInTx(t *testing.T) {
	ctx := context.Background()
	tx := &fakeTx{}
	users := infrastructure.NewMemoryUserRepository()
	clock := test.NewFixedClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	uc := application.NewRegisterUserUseCase(tx, users, clock, fixedID("usr_1"))
	email, _ := domain.NewEmail("taro@example.com")

	id, err := uc.Execute(ctx, application.RegisterUserCommand{Email: email, Name: "Taro"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if id != "usr_1" || tx.calls != 1 {
		t.Errorf("id = %q, tx calls = %d", id, tx.calls)
	}
	got, err := users.FindByEmail(ctx, email)
	if err != nil || got.ID() != "usr_1" {
		t.Errorf("FindByEmail = %v, %v", got, err)
	}
}

func TestRegisterUserUseCase_ReturnsErrors(t *testing.T) {
	ctx := context.Background()
	users := infrastructure.NewMemoryUserRepository()
	clock := test.NewFixedClock(time.Now())
	uc := application.NewRegisterUserUseCase(&fakeTx{}, users, clock, fixedID("usr_1"))
	email, _ := domain.NewEmail("taro@example.com")

	if _, err := uc.Execute(ctx, application.RegisterUserCommand{Email: email}); err == nil {
		t.Error("empty name: want error, got nil")
	}
	if _, err := uc.Execute(ctx, application.RegisterUserCommand{Email: email, Name: "Taro"}); err != nil {
		t.Fatal(err)
	}
	_, err := uc.Execute(ctx, application.RegisterUserCommand{Email: email, Name: "Taro"})
	if !errors.Is(err, infrastructure.ErrDuplicateUser) {
		t.Errorf("same id twice: err = %v, want ErrDuplicateUser", err)
	}
}
