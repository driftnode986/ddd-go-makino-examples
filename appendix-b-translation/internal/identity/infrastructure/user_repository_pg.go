package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/infrastructure/sqlcgen"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/pgxtx"
)

// #region user-repository-pg
type UserRepositoryPG struct {
	queries *sqlcgen.Queries
}

func (r *UserRepositoryPG) Add(ctx context.Context, u *domain.User) error {
	return r.q(ctx).InsertUser(ctx, sqlcgen.InsertUserParams{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		Name:      u.Name(),
		CreatedAt: pgtype.Timestamptz{Time: u.CreatedAt(), Valid: true},
	})
}

func (r *UserRepositoryPG) Save(ctx context.Context, u *domain.User) error {
	n, err := r.q(ctx).UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID: u.ID().String(), Email: u.Email().String(), Name: u.Name(),
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// #endregion

func NewUserRepositoryPG(db sqlcgen.DBTX) *UserRepositoryPG {
	return &UserRepositoryPG{queries: sqlcgen.New(db)}
}

func (r *UserRepositoryPG) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	row, err := r.q(ctx).GetUserByID(ctx, id.String())
	return toUser(row, err)
}

func (r *UserRepositoryPG) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	row, err := r.q(ctx).GetUserByEmail(ctx, email.String())
	return toUser(row, err)
}

// q は RunInTx の中ならそのトランザクションで、外なら接続プールでクエリを流す。
func (r *UserRepositoryPG) q(ctx context.Context) *sqlcgen.Queries {
	if tx, ok := pgxtx.TxFromContext(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func toUser(row sqlcgen.User, err error) (*domain.User, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	email, err := domain.NewEmail(row.Email)
	if err != nil {
		return nil, fmt.Errorf("stored email: %w", err)
	}
	return domain.ReconstructUser(domain.UserID(row.ID), email, row.Name, row.CreatedAt.Time), nil
}

var _ domain.UserRepository = (*UserRepositoryPG)(nil)
