package application_test

import (
	"context"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/infrastructure/memory"
)

func TestRegisterMemberUseCase(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemberRepository()
	uc := application.NewRegisterMemberUseCase(repo)

	out, err := uc.Execute(ctx, application.RegisterMemberInput{Email: "hanako@example.com", Name: "Hanako"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	id, err := domain.MemberIDFrom(out.MemberID)
	if err != nil {
		t.Fatal(err)
	}
	m, err := repo.FindByID(ctx, id)
	if err != nil || m.Email().String() != "hanako@example.com" {
		t.Fatalf("saved member = %v, %v", m, err)
	}

	for _, in := range []application.RegisterMemberInput{
		{Email: "not-an-email", Name: "X"},
		{Email: "x@example.com", Name: ""},
	} {
		if _, err := uc.Execute(ctx, in); err == nil {
			t.Errorf("Execute(%+v): want error, got nil", in)
		}
	}
}
