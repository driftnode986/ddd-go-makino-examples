// internal/booking/infrastructure/memory/booking_repository.go
package memory

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
)

type BookingRepository struct {
	mu       sync.RWMutex
	bookings map[domain.BookingID]*domain.Booking
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{bookings: map[domain.BookingID]*domain.Booking{}}
}

func (r *BookingRepository) Save(_ context.Context, b *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[b.ID()] = b
	return nil
}

func (r *BookingRepository) FindByID(_ context.Context,
	id domain.BookingID) (*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bookings[id]
	if !ok {
		return nil, domain.ErrBookingNotFound
	}
	return b, nil
}

func (r *BookingRepository) ListByMember(_ context.Context,
	memberID domain.MemberID) ([]*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Booking, 0)
	for _, b := range r.bookings {
		if b.MemberID() == memberID {
			out = append(out, b)
		}
	}
	return out, nil
}
