// #region postgres-repo
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/infrastructure/sqlcgen"
)

type PostgresBookingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresBookingRepository(pool *pgxpool.Pool) *PostgresBookingRepository {
	return &PostgresBookingRepository{pool: pool}
}

// #endregion

// #region find-by-id
func (r *PostgresBookingRepository) FindByID(
	ctx context.Context, id domain.BookingID,
) (*domain.Booking, error) {
	q := r.queries(ctx)
	row, err := q.GetBookingByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get booking by id: %w", err)
	}
	return toBookingDomain(bookingRow{
		ID:          row.ID,
		MemberID:    row.MemberID,
		LessonID:    row.LessonID,
		Status:      row.Status,
		ScheduledAt: row.ScheduledAt.Time,
		Version:     int(row.Version),
	})
}

// #endregion

// #region add
func (r *PostgresBookingRepository) Add(
	ctx context.Context, b *domain.Booking,
) error {
	row := toBookingRow(b)
	q := r.queries(ctx)
	if err := q.InsertBooking(ctx, sqlcgen.InsertBookingParams{
		ID:          row.ID,
		MemberID:    row.MemberID,
		LessonID:    row.LessonID,
		Status:      row.Status,
		ScheduledAt: pgtype.Timestamptz{Time: row.ScheduledAt, Valid: true},
		Version:     int32(row.Version),
	}); err != nil {
		return fmt.Errorf("insert booking: %w", err)
	}
	return nil
}

// #endregion

// #region save
func (r *PostgresBookingRepository) Save(
	ctx context.Context, b *domain.Booking,
) error {
	row := toBookingRow(b)
	q := r.queries(ctx)
	rowsAffected, err := q.UpdateBookingWithVersion(
		ctx,
		sqlcgen.UpdateBookingWithVersionParams{
			Status:          row.Status,
			ScheduledAt:     pgtype.Timestamptz{Time: row.ScheduledAt, Valid: true},
			Version:         int32(row.Version),
			ID:              row.ID,
			ExpectedVersion: int32(row.Version - 1),
		},
	)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrVersionConflict
	}
	return nil
}

// #endregion

// Remove は予約を削除する。該当する行が無ければ domain.ErrNotFound を返す
// (InMemoryBookingRepository と同じ振る舞い)。
func (r *PostgresBookingRepository) Remove(
	ctx context.Context, id domain.BookingID,
) error {
	rowsAffected, err := r.queries(ctx).DeleteBooking(ctx, id.String())
	if err != nil {
		return fmt.Errorf("delete booking: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// #region list-by-member
func (r *PostgresBookingRepository) ListByMember(
	ctx context.Context, memberID domain.MemberID,
) ([]*domain.Booking, error) {
	rows, err := r.queries(ctx).ListBookingsByMember(ctx, memberID.String())
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}
	bookings := make([]*domain.Booking, 0, len(rows))
	for _, row := range rows {
		b, err := toBookingDomain(rowToBookingRow(row))
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// #endregion

var _ domain.Repository = (*PostgresBookingRepository)(nil)
