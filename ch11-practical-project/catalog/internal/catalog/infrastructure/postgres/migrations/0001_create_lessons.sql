CREATE SCHEMA catalog;

CREATE TABLE catalog.lessons (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    capacity   INTEGER NOT NULL,
    price_jpy  INTEGER NOT NULL,
    starts_at  TIMESTAMPTZ NOT NULL,
    trainer_id UUID NOT NULL
);
