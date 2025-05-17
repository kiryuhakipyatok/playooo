-- +goose Up
-- +goose StatementBegin
CREATE TABLE games(
    id VARCHAR(45) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(45) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    banner TEXT NOT NULL,
    picture TEXT NOT NULL,
    number_of_players INT DEFAULT 0,
    number_of_events INT DEFAULT 0,
    rating NUMERIC DEFAULT 0.0
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games
-- +goose StatementEnd
