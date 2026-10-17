// #region memory-repo
package infrastructure

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
)

type InMemoryBookingRepository struct {
	mu       sync.Mutex
	bookings map[domain.BookingID]*domain.Booking
}

func NewInMemoryBookingRepository() *InMemoryBookingRepository {
	return &InMemoryBookingRepository{
		bookings: map[domain.BookingID]*domain.Booking{},
	}
}

// #endregion

// #region memory-add-save
func (r *InMemoryBookingRepository) Add(
	ctx context.Context, b *domain.Booking,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bookings[b.ID()]; exists {
		return domain.ErrDuplicate
	}
	r.bookings[b.ID()] = copyBooking(b)
	return nil
}

func (r *InMemoryBookingRepository) Save(
	ctx context.Context, b *domain.Booking,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.bookings[b.ID()]
	if !ok {
		return domain.ErrNotFound
	}
	if existing.Version() != b.Version()-1 {
		return domain.ErrVersionConflict
	}
	r.bookings[b.ID()] = copyBooking(b)
	return nil
}

// #endregion

// #region memory-find-remove
func (r *InMemoryBookingRepository) FindByID(
	ctx context.Context, id domain.BookingID,
) (*domain.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bookings[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return copyBooking(b), nil
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
	return nil
}

// #endregion

// #region copy-booking
// copyBooking は集約のコピーを返す。保存した集約と呼び出し側の集約が
// 同じポインタを共有すると、Confirm などで呼び出し側の version が
// 進んだとき保存済みの version も進み、Save の競合判定が働かない。
// PostgreSQL と同じく「保存した時点の状態」を持つためにコピーする。
func copyBooking(b *domain.Booking) *domain.Booking {
	c := *b
	return &c
}

// #endregion
