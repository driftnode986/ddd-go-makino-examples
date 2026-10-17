package pgxtx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/pgxtx"
)

func startPostgres(t *testing.T, schema string) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
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

func count(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM items").Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestTxManager_CommitsAndRollsBack(t *testing.T) {
	ctx := context.Background()
	pool := startPostgres(t, "CREATE TABLE items (id INT PRIMARY KEY)")
	m := pgxtx.NewTxManager(pool)
	insert := func(id int) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			tx, ok := pgxtx.TxFromContext(ctx)
			if !ok {
				return errors.New("no tx in context")
			}
			_, err := tx.Exec(ctx, "INSERT INTO items (id) VALUES ($1)", id)
			return err
		}
	}

	if err := m.RunInTx(ctx, insert(1)); err != nil {
		t.Fatalf("commit: %v", err)
	}
	errBoom := errors.New("boom")
	err := m.RunInTx(ctx, func(ctx context.Context) error {
		if err := insert(2)(ctx); err != nil {
			return err
		}
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if n := count(t, pool); n != 1 {
		t.Errorf("rows = %d, want 1 (the second insert must be rolled back)", n)
	}
	if _, ok := pgxtx.TxFromContext(ctx); ok {
		t.Error("TxFromContext outside RunInTx: want false")
	}
}
