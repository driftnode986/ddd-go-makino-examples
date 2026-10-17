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
	u, err := domain.NewUser(id, mustEmail(t, "a@example.com"), "Alice")
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if u.ID() != id || u.Email().String() != "a@example.com" || u.DisplayName() != "Alice" {
		t.Fatalf("fields not set: %v %v %q", u.ID(), u.Email(), u.DisplayName())
	}
	if _, err := domain.NewUser(id, mustEmail(t, "a@example.com"), ""); !errors.Is(err, domain.ErrInvalidDisplayName) {
		t.Fatalf("empty name: want ErrInvalidDisplayName, got %v", err)
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	u, err := domain.NewUser(shared.NewUserID(), mustEmail(t, "a@example.com"), "Alice")
	if err != nil {
		t.Fatal(err)
	}
	if err := u.ChangeEmail(mustEmail(t, "a@example.com")); err == nil {
		t.Fatal("same email: want error")
	}
	if err := u.ChangeEmail(mustEmail(t, "b@example.com")); err != nil {
		t.Fatalf("ChangeEmail: %v", err)
	}
	if u.Email().String() != "b@example.com" {
		t.Fatalf("Email: got %v", u.Email())
	}
}

func TestMoney_AddAndEquals(t *testing.T) {
	a, err := domain.NewMoney(1000, "JPY")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := domain.NewMoney(500, "JPY")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := domain.NewMoney(1500, "JPY")
	if !sum.Equals(want) {
		t.Fatalf("Add: got %v want %v", sum, want)
	}
	if !a.Equals(mustMoney(t, 1000, "JPY")) {
		t.Fatal("Add must not change the receiver")
	}
	usd, _ := domain.NewMoney(10, "USD")
	if _, err := a.Add(usd); err == nil {
		t.Fatal("currency mismatch: want error")
	}
	if _, err := domain.NewMoney(-1, "JPY"); err == nil {
		t.Fatal("negative amount: want error")
	}
}

func mustMoney(t *testing.T, amount int64, c domain.Currency) domain.Money {
	t.Helper()
	m, err := domain.NewMoney(amount, c)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
