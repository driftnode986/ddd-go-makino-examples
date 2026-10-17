package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

func (u *BookLessonUseCase) Execute(
	ctx context.Context,
	in BookLessonInput,
) error {
	booking, err := domain.NewBooking(
		in.MemberID,
		in.LessonID,
		in.ScheduledAt,
	)
	if err != nil {
		return fmt.Errorf("create booking: %w", err)
	}

	if err := u.repo.Save(ctx, booking); err != nil {
		var conflict *domain.ScheduleConflictError
		if errors.As(err, &conflict) {
			return fmt.Errorf("conflict with %s: %w",
				conflict.ExistingBookingID, err)
		}
		return fmt.Errorf("save booking: %w", err)
	}
	return nil
}
