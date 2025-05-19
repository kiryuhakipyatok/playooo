-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_notifications(
    user_id UUID NOT NULL,
    notification_id UUID NOT NULL,
    PRIMARY KEY (user_id,notification_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_notifications
-- +goose StatementEnd
