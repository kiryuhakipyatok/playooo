-- +goose Up
-- +goose StatementBegin
CREATE TABLE events_comments(
    comment_id UUID NOT NULL UNIQUE,
    event_id UUID NOT NULL,
    PRIMARY KEY (comment_id,event_id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events_comments
-- +goose StatementEnd
