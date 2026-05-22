package application

import (
	"context"
	"sync"
)

type IdempotencyStore interface {
	AlreadyProcessed(ctx context.Context, eventID int64) (bool, error)
	MarkProcessed(ctx context.Context, eventID int64) error
}

type MemoryIdempotencyStore struct {
	mu        sync.Mutex
	processed map[int64]struct{}
}

func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	return &MemoryIdempotencyStore{processed: map[int64]struct{}{}}
}

func (s *MemoryIdempotencyStore) AlreadyProcessed(
	ctx context.Context, eventID int64,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.processed[eventID]
	return ok, nil
}

func (s *MemoryIdempotencyStore) MarkProcessed(
	ctx context.Context, eventID int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[eventID] = struct{}{}
	return nil
}
