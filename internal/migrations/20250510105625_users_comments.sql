-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_comments(
    comment_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    PRIMARY KEY (comment_id,user_id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_comments
-- +goose StatementEnd
