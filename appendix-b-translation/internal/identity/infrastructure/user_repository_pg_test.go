package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/db"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/pgxtx"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

func startPostgres(t *testing.T, schema string) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog(
				"database system is ready to accept connections",
			).WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("postgres.Run: %v", err)
	}
	t.Cleanup(func() { _ = pg.Terminate(ctx) })
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if _, err := pool.Exec(ctx, schema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}
	return pool
}

func TestUserRepositoryPG(t *testing.T) {
	ctx := context.Background()
	pool := startPostgres(t, db.IdentitySchema)
	repo := infrastructure.NewUserRepositoryPG(pool)
	email, _ := domain.NewEmail("taro@example.com")
	created := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	u, _ := domain.NewUser("usr_1", email, "Taro", created)

	if err := repo.Add(ctx, u); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := u.Rename("Jiro"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, u); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := repo.FindByEmail(ctx, email)
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if got.ID() != "usr_1" || got.Name() != "Jiro" ||
		!got.CreatedAt().Equal(created) {
		t.Errorf(
			"user = %v %v %v",
			got.ID(),
			got.Name(),
			got.CreatedAt(),
		)
	}
	if _, err := repo.FindByID(ctx, "usr_missing"); !errors.Is(
		err,
		domain.ErrUserNotFound,
	) {
		t.Errorf(
			"FindByID(missing) err = %v, want ErrUserNotFound",
			err,
		)
	}
	missing, _ := domain.NewUser("usr_missing", email, "X", created)
	if err := repo.Save(ctx, missing); !errors.Is(
		err,
		domain.ErrUserNotFound,
	) {
		t.Errorf("Save(missing) err = %v, want ErrUserNotFound", err)
	}
}

// RunInTx の中で Add した会員は、fn がエラーを返すとロールバックされて残らない。
func TestRegisterUserUseCase_RollsBackWithPostgres(t *testing.T) {
	ctx := context.Background()
	pool := startPostgres(t, db.IdentitySchema)
	repo := infrastructure.NewUserRepositoryPG(pool)
	txMgr := pgxtx.NewTxManager(pool)
	clock := test.NewFixedClock(
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	uc := application.NewRegisterUserUseCase(
		txMgr,
		repo,
		clock,
		infrastructure.RandomIDGenerator{},
	)
	email, _ := domain.NewEmail("hanako@example.com")

	id, err := uc.Execute(
		ctx,
		application.RegisterUserCommand{
			Email: email,
			Name:  "Hanako",
		},
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if _, err := repo.FindByID(ctx, id); err != nil {
		t.Fatalf("committed user not found: %v", err)
	}

	errBoom := errors.New("boom")
	other, _ := domain.NewEmail("jiro@example.com")
	err = txMgr.RunInTx(ctx, func(ctx context.Context) error {
		u, _ := domain.NewUser(
			"usr_rollback",
			other,
			"Jiro",
			clock.Now(),
		)
		if err := repo.Add(ctx, u); err != nil {
			return err
		}
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("RunInTx err = %v, want boom", err)
	}
	if _, err := repo.FindByID(ctx, "usr_rollback"); !errors.Is(
		err,
		domain.ErrUserNotFound,
	) {
		t.Errorf(
			"rolled back user: err = %v, want ErrUserNotFound",
			err,
		)
	}
}
