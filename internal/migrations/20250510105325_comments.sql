-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments(
    id UUID PRIMARY KEY NOT NULL,
    author_id UUID REFERENCES users(id) NOT NULL,
    author_login VARCHAR(45) NOT NULL,
    author_avatar VARCHAR(45),
    body text NOT NULL,
   -- receiver_id UUID REFERENCES users(id) NOT NULL,
    time TIMESTAMPTZ NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments
-- +goose StatementEnd
