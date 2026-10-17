-- name: GetBookingByID :one
SELECT id, member_id, lesson_id, status, version, created_at
FROM booking.bookings
WHERE id = $1;

-- name: SaveBooking :exec
INSERT INTO booking.bookings (
    id, member_id, lesson_id, status, version, created_at
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE
SET status = EXCLUDED.status, version = EXCLUDED.version;
