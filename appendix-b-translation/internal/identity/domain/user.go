// #region user
package domain

import (
	"errors"
	"time"
)

type User struct {
	id        UserID
	email     Email
	name      string
	createdAt time.Time
}

func NewUser(id UserID, email Email, name string, now time.Time) (*User, error) {
	if name == "" {
		return nil, errors.New("name required")
	}
	return &User{id: id, email: email, name: name, createdAt: now}, nil
}

func (u *User) Rename(newName string) error {
	if newName == "" {
		return errors.New("name required")
	}
	u.name = newName
	return nil
}

// #endregion

// ReconstructUser は保存済みの値から User を組み立て直す（リポジトリが使う）。
func ReconstructUser(id UserID, email Email, name string, createdAt time.Time) *User {
	return &User{id: id, email: email, name: name, createdAt: createdAt}
}

func (u *User) ID() UserID           { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) Name() string         { return u.name }
func (u *User) CreatedAt() time.Time { return u.createdAt }
