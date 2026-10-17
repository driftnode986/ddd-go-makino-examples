package infrastructure_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/db"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
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

func newOrder(t *testing.T) *domain.Order {
	t.Helper()
	o, err := domain.PlaceOrder([]domain.OrderItem{
		{ProductID: "p1", Quantity: 2, UnitPrice: 1500},
		{ProductID: "p2", Quantity: 1, UnitPrice: 800},
	}, "cus_001", time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func TestOrderRepositories(t *testing.T) {
	repos := map[string]domain.OrderRepository{
		"memory":   infra.NewInMemoryOrderRepository(),
		"postgres": infra.NewOrderRepositoryPG(startPostgres(t, db.BookingSchema)),
	}
	for name, repo := range repos {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			o := newOrder(t)
			if err := repo.Add(ctx, o); err != nil {
				t.Fatalf("Add: %v", err)
			}
			if err := repo.Add(ctx, o); err == nil {
				t.Error("second Add: want error, got nil")
			}
			got, err := repo.FindByID(ctx, o.ID())
			if err != nil {
				t.Fatalf("FindByID: %v", err)
			}
			if diff := cmp.Diff(o, got, cmp.AllowUnexported(domain.Order{})); diff != "" {
				t.Errorf("order mismatch (-want +got):\n%s", diff)
			}
			if _, err := repo.FindByID(ctx, "ord_missing"); !errors.Is(err, domain.ErrOrderNotFound) {
				t.Errorf("missing: err = %v, want ErrOrderNotFound", err)
			}
		})
	}
}
