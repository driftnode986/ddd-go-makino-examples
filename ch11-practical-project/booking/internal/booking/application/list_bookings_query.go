package application

import (
	"context"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
)

type BookingView struct {
	ID       string
	MemberID string
	LessonID string
	Status   string
	Version  int
}

type ListBookingsQuery struct {
	repo domain.BookingRepository
}

func NewListBookingsQuery(
	repo domain.BookingRepository,
) *ListBookingsQuery {
	return &ListBookingsQuery{repo: repo}
}

func (q *ListBookingsQuery) ByMember(ctx context.Context,
	memberIDStr string) ([]BookingView, error) {
	memberID, err := domain.MemberIDFrom(memberIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid member id: %w", err)
	}
	bookings, err := q.repo.ListByMember(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	out := make([]BookingView, 0, len(bookings))
	for _, b := range bookings {
		out = append(out, BookingView{
			ID:       b.ID().String(),
			MemberID: b.MemberID().String(),
			LessonID: b.LessonID().String(),
			Status:   string(b.Status()),
			Version:  b.Version(),
		})
	}
	return out, nil
}
