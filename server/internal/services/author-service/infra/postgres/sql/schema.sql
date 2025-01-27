CREATE TABLE authors (
    id         BIGINT PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    user_id    BIGINT NOT NULL,
    metadata   JSONB
);