-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    login VARCHAR(45) UNIQUE NOT NULL,
    telegram VARCHAR(45) UNIQUE NOT NULL,
    chat_id VARCHAR(45) UNIQUE DEFAULT 'unknown',
    rating NUMERIC DEFAULT 0,
    total_rating INT DEFAULT 0,
    number_of_ratings INT DEFAULT 0,
    games text[],
    password BYTEA NOT NULL,
    avatar TEXT DEFAULT 'absent',
    discord VARCHAR(45) UNIQUE DEFAULT 'unknown'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
