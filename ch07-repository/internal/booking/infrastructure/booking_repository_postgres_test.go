// #region postgres-test-imports
//go:build integration

package infrastructure_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure"
)

// #endregion

// #region postgres-test-start
func TestPostgresBookingRepository_AddAndFind(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
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
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	// #endregion

	// #region postgres-test-connect
	connStr, err := pgContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
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

	// #endregion

	// #region postgres-test-assert
	repo := infrastructure.NewPostgresBookingRepository(pool)
	booking := mustNewBooking(t, "2026-06-01T10:00:00Z")
	if err := repo.Add(ctx, booking); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, err := repo.FindByID(ctx, booking.ID())
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if diff := cmp.Diff(
		booking, got,
		cmp.AllowUnexported(
			domain.Booking{}, domain.BookingID{},
			domain.MemberID{}, domain.LessonID{},
		),
	); diff != "" {
		t.Errorf("Booking mismatch (-want +got):\n%s", diff)
	}
}

// #endregion
