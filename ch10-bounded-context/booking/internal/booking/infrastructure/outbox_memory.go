package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/booking/internal/booking/domain"
)

type OutboxEntry struct {
	ID          int64
	EventName   string
	AggregateID string
	Payload     []byte
	OccurredAt  time.Time
	PublishedAt *time.Time
}

type MemoryOutbox struct {
	mu      sync.Mutex
	nextID  int64
	entries []OutboxEntry
}

func NewMemoryOutbox() *MemoryOutbox {
	return &MemoryOutbox{}
}

func (o *MemoryOutbox) Append(
	ctx context.Context, events []domain.DomainEvent,
) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, ev := range events {
		payload, err := marshalEvent(ev)
		if err != nil {
			return fmt.Errorf("marshal %s: %w", ev.EventName(), err)
		}
		o.nextID++
		o.entries = append(o.entries, OutboxEntry{
			ID:          o.nextID,
			EventName:   ev.EventName(),
			AggregateID: ev.AggregateID(),
			Payload:     payload,
			OccurredAt:  ev.OccurredAt(),
		})
	}
	return nil
}

// Unpublished は未送信のイベントを id 順で返す (publisher 用)。
func (o *MemoryOutbox) Unpublished() []OutboxEntry {
	o.mu.Lock()
	defer o.mu.Unlock()
	result := make([]OutboxEntry, 0)
	for _, e := range o.entries {
		if e.PublishedAt == nil {
			result = append(result, e)
		}
	}
	return result
}

// MarkPublished は publisher が送信成功後に呼ぶ。
func (o *MemoryOutbox) MarkPublished(id int64, at time.Time) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for i, e := range o.entries {
		if e.ID == id {
			t := at
			o.entries[i].PublishedAt = &t
			return
		}
	}
}

// marshalEvent はドメインイベントを JSON にシリアライズする。
// イベントごとに json タグを持つ DTO に詰め替えてから marshal する。
func marshalEvent(ev domain.DomainEvent) ([]byte, error) {
	switch e := ev.(type) {
	case domain.OrderPlaced:
		return json.Marshal(struct {
			OrderID     string    `json:"order_id"`
			CustomerID  string    `json:"customer_id"`
			TotalAmount int64     `json:"total_amount"`
			OccurredAt  time.Time `json:"occurred_at"`
		}{
			OrderID:     e.OrderID().String(),
			CustomerID:  e.CustomerID().String(),
			TotalAmount: e.TotalAmount(),
			OccurredAt:  e.OccurredAt(),
		})
	default:
		return nil, fmt.Errorf("unknown event: %s", ev.EventName())
	}
}
