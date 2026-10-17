// Package pgxtx は pgx のトランザクションを context に載せて渡す TxManager を置く。
// RunInTx の中で呼ばれたリポジトリは、TxFromContext(ctx) で同じトランザクションを使う。
package pgxtx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(
	pool *pgxpool.Pool,
) *TxManager {
	return &TxManager{pool: pool}
}

// RunInTx は fn を 1 つのトランザクションで実行する。fn がエラーを返したらロールバックする。
func (m *TxManager) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("%w (rollback: %v)", err, rbErr)
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// TxFromContext は RunInTx が ctx に載せたトランザクションを返す。
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
