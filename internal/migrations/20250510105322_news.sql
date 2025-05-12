-- +goose Up
-- +goose StatementBegin
CREATE TABLE news(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    title VARCHAR(45) NOT NULL,
    body TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    link TEXT DEFAULT 'absent',
    picture TEXT DEFAULT 'absent'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE news
-- +goose StatementEnd
