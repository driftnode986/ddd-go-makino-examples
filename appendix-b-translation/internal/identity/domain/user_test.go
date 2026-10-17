package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
)

func TestNewEmail_NormalizesDomain(t *testing.T) {
	e, err := domain.NewEmail("Taro@Example.COM")
	if err != nil {
		t.Fatal(err)
	}
	if e.String() != "Taro@example.com" {
		t.Errorf("email = %q, want Taro@example.com", e.String())
	}
	if _, err := domain.NewEmail("not an email"); err == nil {
		t.Error("NewEmail(invalid): want error, got nil")
	}
}

func TestNewUserAndRename(t *testing.T) {
	email, _ := domain.NewEmail("taro@example.com")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	u, err := domain.NewUser("usr_1", email, "Taro", now)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID() != "usr_1" || u.Name() != "Taro" ||
		!u.CreatedAt().Equal(now) {
		t.Errorf("user = %+v", u)
	}
	if err := u.Rename("Jiro"); err != nil || u.Name() != "Jiro" {
		t.Errorf("Rename: name = %q, err = %v", u.Name(), err)
	}
	if err := u.Rename(""); err == nil {
		t.Error("Rename(\"\"): want error, got nil")
	}
	if _, err := domain.NewUser(
		"usr_2",
		email,
		"",
		now,
	); err == nil {
		t.Error("NewUser with empty name: want error, got nil")
	}
}
