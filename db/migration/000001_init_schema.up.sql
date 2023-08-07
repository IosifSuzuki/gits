CREATE TABLE IF NOT EXISTS account
(
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE,
    hash_password TEXT,
    role VARCHAR(32),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)