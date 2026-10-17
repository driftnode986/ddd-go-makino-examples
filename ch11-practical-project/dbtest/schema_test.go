package dbtest_test

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// 各 BC のマイグレーション。3 BC が 1 つのデータベースに同居する構成を作る。
var migrations = []string{
	"../booking/internal/booking/infrastructure/postgres/migrations",
	"../catalog/internal/catalog/infrastructure/postgres/migrations",
	"../identity/internal/identity/" +
		"infrastructure/postgres/migrations",
}

const bookingQueries = "../booking/internal/booking/" +
	"infrastructure/postgres/queries.sql"

func startPostgres(t *testing.T) *pgx.Conn {
	t.Helper()
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:17-alpine",
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

	for _, dir := range migrations {
		files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
		if err != nil || len(files) == 0 {
			t.Fatalf("no migrations in %s: %v", dir, err)
		}
		for _, f := range files {
			sql, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := conn.Exec(ctx, string(sql)); err != nil {
				t.Fatalf("apply %s: %v", f, err)
			}
		}
	}
	return conn
}

// loadQueries は sqlc 形式（-- name: Xxx :kind）の queries.sql を名前ごとに分ける。
func loadQueries(t *testing.T, path string) map[string]string {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`(?m)^-- name: (\w+) :\w+\n`)
	idx := re.FindAllStringSubmatchIndex(string(src), -1)
	out := map[string]string{}
	for i, m := range idx {
		end := len(src)
		if i+1 < len(idx) {
			end = idx[i+1][0]
		}
		out[string(src[m[2]:m[3]])] = strings.TrimSpace(
			string(src[m[1]:end]),
		)
	}
	return out
}

func TestMigrations_KeepEachBCInItsOwnSchemaWithoutForeignKeys(
	t *testing.T,
) {
	conn := startPostgres(t)
	ctx := context.Background()

	for _, table := range []string{
		"booking.bookings",
		"catalog.lessons",
		"identity.members",
	} {
		var exists bool
		if err := conn.QueryRow(
			ctx,
			`SELECT to_regclass($1) IS NOT NULL`,
			table,
		).Scan(&exists); err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Errorf("table %s was not created", table)
		}
	}
	var fks int
	err := conn.QueryRow(ctx, `SELECT count(*)
		FROM information_schema.table_constraints
		WHERE constraint_type = 'FOREIGN KEY'
		AND table_schema IN ('booking', 'catalog', 'identity')`).
		Scan(&fks)
	if err != nil {
		t.Fatal(err)
	}
	if fks != 0 {
		t.Errorf("foreign keys between BCs = %d, want 0", fks)
	}
}

func TestBookingQueries_SaveAndGet(t *testing.T) {
	conn := startPostgres(t)
	ctx := context.Background()
	q := loadQueries(t, bookingQueries)
	if len(q) != 2 {
		t.Fatalf(
			"queries = %v, want GetBookingByID and SaveBooking",
			q,
		)
	}

	id, member, lesson := uuid.New(), uuid.New(), uuid.New()
	created := time.Date(2026, 5, 23, 10, 0, 0, 0, time.UTC)
	if _, err := conn.Exec(
		ctx,
		q["SaveBooking"],
		id,
		member,
		lesson,
		"pending",
		1,
		created,
	); err != nil {
		t.Fatalf("SaveBooking (insert): %v", err)
	}
	// 同じ id で保存し直すと status と version だけが更新される
	// （ON CONFLICT ... DO UPDATE）
	if _, err := conn.Exec(
		ctx,
		q["SaveBooking"],
		id,
		member,
		lesson,
		"confirmed",
		2,
		created,
	); err != nil {
		t.Fatalf("SaveBooking (update): %v", err)
	}

	var (
		gotID, gotMember, gotLesson uuid.UUID
		status                      string
		version                     int64
		createdAt                   time.Time
	)
	err := conn.QueryRow(ctx, q["GetBookingByID"], id).
		Scan(
			&gotID,
			&gotMember,
			&gotLesson,
			&status,
			&version,
			&createdAt,
		)
	if err != nil {
		t.Fatalf("GetBookingByID: %v", err)
	}
	if gotID != id || gotMember != member || gotLesson != lesson {
		t.Errorf("ids = %v %v %v", gotID, gotMember, gotLesson)
	}
	if status != "confirmed" || version != 2 ||
		!createdAt.Equal(created) {
		t.Errorf(
			"status = %s, version = %d, created_at = %v",
			status,
			version,
			createdAt,
		)
	}
}
