CREATE TABLE orders (
    id          TEXT PRIMARY KEY,
    customer_id TEXT NOT NULL,
    status      TEXT NOT NULL,
    items       JSONB NOT NULL,
    placed_at   TIMESTAMPTZ NOT NULL
);
