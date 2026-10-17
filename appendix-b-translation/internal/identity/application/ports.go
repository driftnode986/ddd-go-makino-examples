package application

import (
	"context"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
)

// TxManager は fn を 1 つのトランザクションの中で実行する。
// fn がエラーを返したらロールバックし、nil ならコミットする。
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// Clock は現在時刻を返す。テストでは固定の時刻に差し替える。
type Clock interface {
	Now() time.Time
}

// IDGenerator は新しい会員 ID を作る。
type IDGenerator interface {
	New() domain.UserID
}

// RegisterUserCommand は RegisterUserUseCase の入力。
type RegisterUserCommand struct {
	Email domain.Email
	Name  string
}
