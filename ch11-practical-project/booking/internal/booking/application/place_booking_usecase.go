// internal/booking/application/place_booking_usecase.go
package application

import (
	"context"
	"fmt"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

// LessonQueryPort は Catalog BC への参照ポートです。
// Booking BC は catalog/internal/... を import せず、shared/dto を介して連携します。
type LessonQueryPort interface {
	FindLessonView(ctx context.Context, lessonID string) (dto.LessonView, error)
}

type PlaceBookingInput struct {
	MemberID string
	LessonID string
}

type PlaceBookingOutput struct {
	BookingID string
	Status    string
}

type PlaceBookingUseCase struct {
	repo        domain.BookingRepository
	lessonQuery LessonQueryPort
	now         func() time.Time
}

func NewPlaceBookingUseCase(repo domain.BookingRepository,
	lessonQuery LessonQueryPort, now func() time.Time) *PlaceBookingUseCase {
	return &PlaceBookingUseCase{repo: repo, lessonQuery: lessonQuery, now: now}
}

func (u *PlaceBookingUseCase) Execute(ctx context.Context,
	in PlaceBookingInput) (PlaceBookingOutput, error) {
	memberID, err := domain.MemberIDFrom(in.MemberID)
	if err != nil {
		return PlaceBookingOutput{}, fmt.Errorf("invalid input: %w", err)
	}
	lessonID, err := domain.LessonIDFrom(in.LessonID)
	if err != nil {
		return PlaceBookingOutput{}, fmt.Errorf("invalid input: %w", err)
	}

	if _, err := u.lessonQuery.FindLessonView(ctx, in.LessonID); err != nil {
		return PlaceBookingOutput{}, fmt.Errorf("lesson not available: %w", err)
	}

	b := domain.NewBooking(memberID, lessonID, u.now())
	if err := u.repo.Save(ctx, b); err != nil {
		return PlaceBookingOutput{}, fmt.Errorf("save: %w", err)
	}
	return PlaceBookingOutput{
		BookingID: b.ID().String(),
		Status:    string(b.Status()),
	}, nil
}
