// #region tx-a
package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure/sqlcgen"
)

type txKey struct{}

// WithTx は context にトランザクションを乗せる。
// ユースケース層が pool.BeginTx で開始したトランザクションを
// ここで運ぶ。
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// #endregion
// #region tx-b
// queries は context にトランザクションが乗っていれば tx、
// なければ pool を使う。
func (r *PostgresBookingRepository) queries(
	ctx context.Context,
) *sqlcgen.Queries {
	if tx, ok := extractTx(ctx); ok {
		return sqlcgen.New(tx)
	}
	return sqlcgen.New(r.pool)
}

// #endregion
