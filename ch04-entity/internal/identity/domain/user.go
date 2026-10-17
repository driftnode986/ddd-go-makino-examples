package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// UserID は User エンティティの ID。
type UserID struct {
	value string
}

// NewUserID は ID 文字列を検証して UserID を返す。
func NewUserID(s string) (UserID, error) {
	if len(s) < 8 {
		return UserID{}, errors.New("user id too short")
	}
	return UserID{value: s}, nil
}

// String は ID 文字列を返す。
func (u UserID) String() string { return u.value }

// IsZero はゼロ値かどうかを判定する。
func (u UserID) IsZero() bool { return u.value == "" }

// UserStatus は User の状態を表す。
type UserStatus string

const (
	// UserActive はアクティブな利用者の状態。
	UserActive UserStatus = "active"
	// UserDeactivated は退会済みの状態。
	UserDeactivated UserStatus = "deactivated"
)

// User は利用者を表すエンティティ。ID で識別される。
type User struct {
	id        UserID
	email     string
	name      string
	status    UserStatus
	version   int64
	createdAt time.Time
	updatedAt time.Time
}

var (
	// ErrAlreadyDeactivated はすでに退会済みの利用者を表す。
	ErrAlreadyDeactivated = errors.New("user already deactivated")
	// ErrCannotUpdateDeactivated は退会済み利用者の更新を表す。
	ErrCannotUpdateDeactivated = errors.New(
		"cannot update deactivated user",
	)
	// ErrInvalidName は不正な氏名を表す。
	ErrInvalidName = errors.New("invalid name")
)

// NewUser は User エンティティを生成する。
func NewUser(
	id UserID,
	email, name string,
	clock func() time.Time,
) (*User, error) {
	if id.IsZero() {
		return nil, errors.New("user id required")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf(
			"%w: email must contain @",
			ErrInvalidName,
		)
	}
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf(
			"%w: name must not be empty",
			ErrInvalidName,
		)
	}
	now := clock()
	return &User{
		id:        id,
		email:     email,
		name:      strings.TrimSpace(name),
		status:    UserActive,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Rename は氏名を変更する。Active のみ。
func (u *User) Rename(name string, clock func() time.Time) error {
	if u.status != UserActive {
		return ErrCannotUpdateDeactivated
	}
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf(
			"%w: name must not be empty",
			ErrInvalidName,
		)
	}
	u.name = trimmed
	u.version++
	u.updatedAt = clock()
	return nil
}

// Deactivate は退会処理を行う。すでに退会済みの場合はエラー。
func (u *User) Deactivate(clock func() time.Time) error {
	if u.status == UserDeactivated {
		return ErrAlreadyDeactivated
	}
	u.status = UserDeactivated
	u.version++
	u.updatedAt = clock()
	return nil
}

// 状態取得。

// ID は User ID を返す。
func (u User) ID() UserID { return u.id }

// Email はメールアドレスを返す。
func (u User) Email() string { return u.email }

// Name は氏名を返す。
func (u User) Name() string { return u.name }

// Status は状態を返す。
func (u User) Status() UserStatus { return u.status }

// Version は楽観的ロック用の版番号を返す。
func (u User) Version() int64 { return u.version }

// CreatedAt は作成日時を返す。
func (u User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt は最終更新日時を返す。
func (u User) UpdatedAt() time.Time { return u.updatedAt }
