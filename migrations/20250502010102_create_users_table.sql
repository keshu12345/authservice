-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    api_key TEXT UNIQUE,
    is_revoke BOOLEAN NOT NULL DEFAULT FALSE 
);

-- +goose Down
DROP TABLE IF EXISTS users;

