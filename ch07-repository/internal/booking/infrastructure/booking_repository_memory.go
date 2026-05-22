package infrastructure

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
)

// InMemoryBookingRepository は domain.Repository のテストダブル実装。
// 楽観的ロックの振る舞いを再現するため、保存時はバージョンスナップショットを記録する。
type InMemoryBookingRepository struct {
	mu       sync.Mutex
	bookings map[domain.BookingID]*domain.Booking
	versions map[domain.BookingID]int
}

func NewInMemoryBookingRepository() *InMemoryBookingRepository {
	return &InMemoryBookingRepository{
		bookings: map[domain.BookingID]*domain.Booking{},
		versions: map[domain.BookingID]int{},
	}
}

func (r *InMemoryBookingRepository) Add(ctx context.Context, b *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bookings[b.ID()]; exists {
		return domain.ErrDuplicate
	}
	r.bookings[b.ID()] = b
	r.versions[b.ID()] = b.Version()
	return nil
}

func (r *InMemoryBookingRepository) Save(ctx context.Context, b *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	storedVersion, ok := r.versions[b.ID()]
	if !ok {
		return domain.ErrNotFound
	}
	// 楽観的ロック: 保存済みの version が「渡された集約の version - 1」と一致する必要がある。
	if storedVersion != b.Version()-1 {
		return domain.ErrVersionConflict
	}
	r.bookings[b.ID()] = b
	r.versions[b.ID()] = b.Version()
	return nil
}

func (r *InMemoryBookingRepository) FindByID(
	ctx context.Context, id domain.BookingID,
) (*domain.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bookings[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return b, nil
}

func (r *InMemoryBookingRepository) Remove(
	ctx context.Context, id domain.BookingID,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.bookings[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.bookings, id)
	delete(r.versions, id)
	return nil
}

// Interface assertion (compile-time check)
var _ domain.Repository = (*InMemoryBookingRepository)(nil)
