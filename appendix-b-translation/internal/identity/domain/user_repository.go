// #region user-repository
package domain

import "context"

type UserRepository interface {
	Add(ctx context.Context, u *User) error
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
}

// #endregion
