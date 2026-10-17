package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
)

// OutboxPort はアプリケーション層から outbox に書き込む契約。
// 実装は infrastructure 層で in-memory または
// PostgreSQL 版を用意する。
type OutboxPort interface {
	Append(ctx context.Context, events []domain.DomainEvent) error
}
