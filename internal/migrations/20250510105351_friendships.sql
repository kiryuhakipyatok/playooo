-- +goose Up
-- +goose StatementBegin
CREATE TABLE friendships(
    user_id1 UUID NOT NULL,
    user_id2 UUID NOT NULL,
    PRIMARY KEY (user_id1,user_id2),
    FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friendships
-- +goose StatementEnd
