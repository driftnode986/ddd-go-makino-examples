-- name: ListBookingsByMember :many
SELECT id, member_id, lesson_id, status, scheduled_at, version
FROM bookings
WHERE member_id = $1
ORDER BY scheduled_at;
