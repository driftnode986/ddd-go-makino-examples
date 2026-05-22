package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch04-entity/internal/identity/domain"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNewUser_RejectsInvalidName(t *testing.T) {
	id, _ := domain.NewUserID("user-0001")
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)

	_, err := domain.NewUser(id, "alice@example.com", "   ", clock)
	if !errors.Is(err, domain.ErrInvalidName) {
		t.Fatalf("want ErrInvalidName, got %v", err)
	}
}

func TestUser_Rename(t *testing.T) {
	id, _ := domain.NewUserID("user-0001")
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	u, err := domain.NewUser(id, "alice@example.com", "Alice", clock)
	if err != nil {
		t.Fatal(err)
	}

	later := fixedClock(now.Add(1 * time.Hour))
	if err := u.Rename("Alice Bob", later); err != nil {
		t.Fatalf("Rename: %v", err)
	}
	if u.Name() != "Alice Bob" {
		t.Errorf("want Alice Bob, got %q", u.Name())
	}
	if u.Version() != 2 {
		t.Errorf("version: want 2, got %d", u.Version())
	}
	if !u.UpdatedAt().After(u.CreatedAt()) {
		t.Errorf("updatedAt should be after createdAt")
	}
}

func TestUser_Deactivate(t *testing.T) {
	id, _ := domain.NewUserID("user-0001")
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)
	u, err := domain.NewUser(id, "alice@example.com", "Alice", clock)
	if err != nil {
		t.Fatal(err)
	}

	if err := u.Deactivate(clock); err != nil {
		t.Fatalf("Deactivate: %v", err)
	}
	if u.Status() != domain.UserDeactivated {
		t.Errorf("want deactivated, got %v", u.Status())
	}

	// 二度目の退会はエラー
	if err := u.Deactivate(clock); !errors.Is(err, domain.ErrAlreadyDeactivated) {
		t.Errorf("want ErrAlreadyDeactivated, got %v", err)
	}

	// 退会後の Rename は不可
	if err := u.Rename("New Name", clock); !errors.Is(err, domain.ErrCannotUpdateDeactivated) {
		t.Errorf("want ErrCannotUpdateDeactivated, got %v", err)
	}
}

func TestUser_IdentityIsByID(t *testing.T) {
	id, _ := domain.NewUserID("user-0001")
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	clock := fixedClock(now)

	a, _ := domain.NewUser(id, "alice@example.com", "Alice", clock)
	b, _ := domain.NewUser(id, "alice@example.com", "Alice", clock)

	// 同じ ID なら同一エンティティ
	if a.ID().String() != b.ID().String() {
		t.Errorf("same ID expected")
	}

	// 片方の状態を変えても、ID は変わらない
	later := fixedClock(now.Add(1 * time.Hour))
	if err := a.Rename("Renamed", later); err != nil {
		t.Fatal(err)
	}
	if a.ID().String() != b.ID().String() {
		t.Errorf("ID must be preserved across mutations")
	}
}
