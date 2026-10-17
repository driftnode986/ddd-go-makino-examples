// #region user-type
package domain

import (
	"errors"
	"time"
)

// User はエンティティ。フィールドは非公開、ポインタレシーバで統一
type User struct {
	id          UserID
	email       Email
	displayName string
	createdAt   time.Time
}

var (
	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrSameEmail          = errors.New("email unchanged")
)

// #endregion

// #region user-methods
// NewUser はファクトリ関数。不変条件を 1 箇所で検証
func NewUser(
	id UserID,
	email Email,
	displayName string,
) (*User, error) {
	if displayName == "" {
		return nil, ErrInvalidDisplayName
	}
	return &User{
		id:          id,
		email:       email,
		displayName: displayName,
		createdAt:   time.Now(),
	}, nil
}

// ビヘイビア指向の状態変更メソッド (ポインタレシーバ)
func (u *User) ChangeEmail(newEmail Email) error {
	if u.email.Equals(newEmail) {
		return ErrSameEmail
	}
	u.email = newEmail
	return nil
}

// 読み取り専用メソッドもポインタレシーバで一貫
func (u *User) ID() UserID          { return u.id }
func (u *User) Email() Email        { return u.email }
func (u *User) DisplayName() string { return u.displayName }

// #endregion
