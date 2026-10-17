package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/db"
)

// MigrateSchema は db/schema.sql を実行して bookings テーブルを作る。
// 統合テストで、起動したばかりの空の PostgreSQL に対して呼ぶ。
func MigrateSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, db.Schema); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}
	return nil
}
