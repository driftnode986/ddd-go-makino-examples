package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/db"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/infrastructure"
)

func startPostgres(t *testing.T) *pgx.Conn {
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
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close(ctx) })
	return conn
}

// in-memory の outbox が作った行を、本のスキーマの outbox_events に入れ、
// publisher と同じ「published_at IS NULL を id 順に読む →
// published_at を立てる」を確かめる。
func TestOutboxSchema_StoresAndPublishesEvents(t *testing.T) {
	ctx := context.Background()
	conn := startPostgres(t)
	if _, err := conn.Exec(ctx, db.OutboxSchema); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	customerID, _ := domain.NewCustomerID("cus_001")
	o := domain.NewOrder(customerID, time.Unix(1700000000, 0))
	pid, _ := domain.NewProductID("prod_001")
	price, _ := domain.NewMoney(1500, "JPY")
	if err := o.AddItem(pid, 2, price); err != nil {
		t.Fatal(err)
	}
	if err := o.Confirm(time.Unix(1700000100, 0)); err != nil {
		t.Fatal(err)
	}
	outbox := infrastructure.NewMemoryOutbox()
	if err := outbox.Append(ctx, o.PullEvents()); err != nil {
		t.Fatal(err)
	}
	for _, e := range outbox.Unpublished() {
		_, err := conn.Exec(ctx, `INSERT INTO outbox_events
			(event_name, aggregate_id, payload, occurred_at)
			VALUES ($1, $2, $3, $4)`,
			e.EventName, e.AggregateID, e.Payload, e.OccurredAt)
		if err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	var (
		id    int64
		name  string
		total int64
	)
	err := conn.QueryRow(ctx, `SELECT id, event_name,
		(payload->>'total_amount')::bigint
		FROM outbox_events WHERE published_at IS NULL
		ORDER BY id LIMIT 1`).
		Scan(&id, &name, &total)
	if err != nil {
		t.Fatalf("select unpublished: %v", err)
	}
	if name != "booking.OrderPlaced" || total != 3000 {
		t.Errorf("event = %s, total = %d", name, total)
	}

	if _, err := conn.Exec(
		ctx,
		`UPDATE outbox_events SET published_at = now()
		WHERE id = $1`,
		id,
	); err != nil {
		t.Fatal(err)
	}
	var left int
	if err := conn.QueryRow(
		ctx,
		`SELECT count(*) FROM outbox_events
		WHERE published_at IS NULL`,
	).Scan(&left); err != nil {
		t.Fatal(err)
	}
	if left != 0 {
		t.Errorf("unpublished rows = %d, want 0", left)
	}
}
