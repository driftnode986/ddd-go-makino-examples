// #region book-lesson-imports
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// #endregion

// #region book-lesson-execute
func (u *BookLessonUseCase) Execute(ctx context.Context, in BookLessonInput) error {
	booking, err := domain.NewBooking(in.MemberID, in.LessonID, in.ScheduledAt)
	if err != nil {
		return fmt.Errorf("create booking: %w", err)
	}

	if err := u.repo.Save(ctx, booking); err != nil {
		slog.ErrorContext(ctx, "save booking failed",
			"member_id", in.MemberID,
			"lesson_id", in.LessonID,
			"error", err,
		)
		return err
	}

	slog.InfoContext(ctx, "booking placed",
		"booking_id", booking.ID(),
		"member_id", in.MemberID,
		"lesson_id", in.LessonID,
	)
	return nil
}

// #endregion
