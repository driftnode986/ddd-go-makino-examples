CREATE SCHEMA identity;

CREATE TABLE identity.members (
    id    UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name  TEXT NOT NULL
);
