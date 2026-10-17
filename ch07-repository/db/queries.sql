-- name: GetBookingByID :one
SELECT id, member_id, lesson_id, status, scheduled_at, version
FROM bookings
WHERE id = $1;

-- name: InsertBooking :exec
INSERT INTO bookings (
    id, member_id, lesson_id, status, scheduled_at, version
)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateBookingWithVersion :execrows
UPDATE bookings
SET status = @status,
    scheduled_at = @scheduled_at,
    version = @version
WHERE id = @id AND version = @expected_version;

-- name: DeleteBooking :execrows
DELETE FROM bookings WHERE id = $1;
