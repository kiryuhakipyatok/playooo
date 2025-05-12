-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    author_id UUID UNIQUE NOT NULL,
    body TEXT DEFAULT 'absent',
    game VARCHAR(45) NOT NULL,
    max INT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    notificated_pre BOOLEAN NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events
-- +goose StatementEnd
