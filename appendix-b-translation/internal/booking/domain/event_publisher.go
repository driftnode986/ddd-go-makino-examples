// #region event-publisher
package domain

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// #endregion
