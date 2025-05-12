-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL UNIQUE,
    body TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications
-- +goose StatementEnd
