//go:build integration

package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure"
)

// 本に載せた AddAndFind のほかの振る舞いを、1 つのコンテナで t.Run に分けて確かめる。
// 各ケースは bookings テーブルを空にしてから始める。
func TestPostgresBookingRepository_Cases(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}
	ctx := context.Background()
	pool := startPostgres(ctx, t)
	repo := infrastructure.NewPostgresBookingRepository(pool)

	run := func(name string, f func(t *testing.T)) {
		t.Run(name, func(t *testing.T) {
			if _, err := pool.Exec(ctx, "TRUNCATE bookings"); err != nil {
				t.Fatalf("truncate: %v", err)
			}
			f(t)
		})
	}

	run("FindByID_NotFound", func(t *testing.T) {
		_, err := repo.FindByID(ctx, mustBookingID(t, "missing"))
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("got %v, want ErrNotFound", err)
		}
	})

	run("Save_AfterFindByID", func(t *testing.T) {
		if err := repo.Add(ctx, mustNewBooking(t, "2026-06-01T10:00:00Z")); err != nil {
			t.Fatalf("Add: %v", err)
		}
		b, err := repo.FindByID(ctx, mustBookingID(t, "bk_001"))
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if err := b.Confirm(); err != nil { // version 1 → 2
			t.Fatalf("Confirm: %v", err)
		}
		if err := repo.Save(ctx, b); err != nil {
			t.Fatalf("Save: %v", err)
		}
		got, err := repo.FindByID(ctx, b.ID())
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if got.Status() != domain.StatusConfirmed || got.Version() != 2 {
			t.Errorf("got status=%v version=%d, want confirmed/2", got.Status(), got.Version())
		}
	})

	run("Save_VersionConflict", func(t *testing.T) {
		b := mustNewBooking(t, "2026-06-01T10:00:00Z")
		if err := repo.Add(ctx, b); err != nil {
			t.Fatalf("Add: %v", err)
		}
		// DB の version は 1。version 5 の集約を保存すると WHERE version = 4 に一致しない
		stale := mustNewBookingWithVersion(t, b.ID(), 5)
		if err := repo.Save(ctx, stale); !errors.Is(err, domain.ErrVersionConflict) {
			t.Errorf("got %v, want ErrVersionConflict", err)
		}
	})

	run("Remove", func(t *testing.T) {
		b := mustNewBooking(t, "2026-06-01T10:00:00Z")
		if err := repo.Add(ctx, b); err != nil {
			t.Fatalf("Add: %v", err)
		}
		if err := repo.Remove(ctx, b.ID()); err != nil {
			t.Fatalf("Remove: %v", err)
		}
		if _, err := repo.FindByID(ctx, b.ID()); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("after Remove: got %v, want ErrNotFound", err)
		}
		if err := repo.Remove(ctx, b.ID()); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Remove twice: got %v, want ErrNotFound", err)
		}
	})

	run("ListByMember_OrderedByScheduledAt", func(t *testing.T) {
		for _, b := range []*domain.Booking{
			mustNewBookingWithID(t, "bk_002", "mem_001", "2026-06-03T10:00:00Z"),
			mustNewBookingWithID(t, "bk_001", "mem_001", "2026-06-01T10:00:00Z"),
			mustNewBookingWithID(t, "bk_003", "mem_999", "2026-06-02T10:00:00Z"),
		} {
			if err := repo.Add(ctx, b); err != nil {
				t.Fatalf("Add %s: %v", b.ID(), err)
			}
		}
		got, err := repo.ListByMember(ctx, mustMemberID(t, "mem_001"))
		if err != nil {
			t.Fatalf("ListByMember: %v", err)
		}
		if len(got) != 2 || got[0].ID().String() != "bk_001" || got[1].ID().String() != "bk_002" {
			ids := make([]string, 0, len(got))
			for _, b := range got {
				ids = append(ids, b.ID().String())
			}
			t.Errorf("got %v, want [bk_001 bk_002]", ids)
		}
	})

	run("WithTx_Rollback", func(t *testing.T) {
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			t.Fatalf("begin tx: %v", err)
		}
		txCtx := infrastructure.WithTx(ctx, tx)
		b := mustNewBooking(t, "2026-06-01T10:00:00Z")
		if err := repo.Add(txCtx, b); err != nil {
			t.Fatalf("Add in tx: %v", err)
		}
		// 同じトランザクションの中からは見える
		if _, err := repo.FindByID(txCtx, b.ID()); err != nil {
			t.Fatalf("FindByID in tx: %v", err)
		}
		if err := tx.Rollback(ctx); err != nil {
			t.Fatalf("rollback: %v", err)
		}
		// ロールバックしたので pool からは見えない
		if _, err := repo.FindByID(ctx, b.ID()); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("after rollback: got %v, want ErrNotFound", err)
		}
	})

	run("WithTx_Commit", func(t *testing.T) {
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			t.Fatalf("begin tx: %v", err)
		}
		defer func() { _ = tx.Rollback(ctx) }()
		b := mustNewBooking(t, "2026-06-01T10:00:00Z")
		if err := repo.Add(infrastructure.WithTx(ctx, tx), b); err != nil {
			t.Fatalf("Add in tx: %v", err)
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatalf("commit: %v", err)
		}
		if _, err := repo.FindByID(ctx, b.ID()); err != nil {
			t.Errorf("after commit: FindByID: %v", err)
		}
	})
}

// startPostgres は本の AddAndFind と同じ手順で PostgreSQL を起動し、
// スキーマを流した接続プールを返す。
func startPostgres(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
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
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("ConnectionString: %v", err)
	}
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := infrastructure.MigrateSchema(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return pool
}
