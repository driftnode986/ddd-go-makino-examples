package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/identity/domain"
	shared "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/shared/domain"
)

func mustEmail(t *testing.T, s string) domain.Email {
	t.Helper()
	e, err := domain.NewEmail(s)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestNewUser(t *testing.T) {
	id := shared.NewUserID()
	u, err := domain.NewUser(
		id,
		mustEmail(t, "a@example.com"),
		"Alice",
	)
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if u.ID() != id || u.Email().String() != "a@example.com" ||
		u.DisplayName() != "Alice" {
		t.Fatalf(
			"fields not set: %v %v %q",
			u.ID(),
			u.Email(),
			u.DisplayName(),
		)
	}
	if _, err := domain.NewUser(
		id,
		mustEmail(t, "a@example.com"),
		"",
	); !errors.Is(
		err,
		domain.ErrInvalidDisplayName,
	) {
		t.Fatalf(
			"empty name: want ErrInvalidDisplayName, got %v",
			err,
		)
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	u, err := domain.NewUser(
		shared.NewUserID(),
		mustEmail(t, "a@example.com"),
		"Alice",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := u.ChangeEmail(
		mustEmail(t, "a@example.com"),
	); !errors.Is(err, domain.ErrSameEmail) {
		t.Fatalf("same email: got %v, want ErrSameEmail", err)
	}
	if err := u.ChangeEmail(
		mustEmail(t, "b@example.com"),
	); err != nil {
		t.Fatalf("ChangeEmail: %v", err)
	}
	if u.Email().String() != "b@example.com" {
		t.Fatalf("Email: got %v", u.Email())
	}
}
