CREATE SCHEMA booking;

CREATE TABLE booking.bookings (
    id         UUID PRIMARY KEY,
    -- identity.members.id 参照だが FK は張らない
    member_id  UUID NOT NULL,
    -- catalog.lessons.id 参照だが FK は張らない
    lesson_id  UUID NOT NULL,
    status     TEXT NOT NULL,
    version    BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
