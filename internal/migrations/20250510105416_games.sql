-- +goose Up
-- +goose StatementBegin
CREATE TABLE games(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(45) NOT NULL UNIQUE,
    number_of_players INT DEFAULT 0,
    number_of_events INT DEFAULT 0,
    rating NUMERIC DEFAULT 0.0
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games
-- +goose StatementEnd
