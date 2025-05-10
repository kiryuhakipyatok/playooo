-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_events(
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    PRIMARY KEY (user_id,event_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_events
-- +goose StatementEnd
