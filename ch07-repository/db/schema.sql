CREATE TABLE bookings (
    id           TEXT PRIMARY KEY,
    member_id    TEXT NOT NULL,
    lesson_id    TEXT NOT NULL,
    status       TEXT NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    version      INTEGER NOT NULL DEFAULT 0
);
