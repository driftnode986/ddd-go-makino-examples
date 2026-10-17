CREATE TABLE outbox_events (
    id           BIGSERIAL PRIMARY KEY,
    event_name   TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    payload      JSONB NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ
);
CREATE INDEX idx_outbox_unpublished
    ON outbox_events (id)
    WHERE published_at IS NULL;
