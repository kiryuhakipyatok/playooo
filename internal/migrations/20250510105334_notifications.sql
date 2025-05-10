-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications(
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES users (id),
    event_id UUID REFERENCES events (id),
    body TEXT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications
-- +goose StatementEnd
