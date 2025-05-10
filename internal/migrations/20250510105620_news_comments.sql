-- +goose Up
-- +goose StatementBegin
CREATE TABLE news_comments(
    comment_id UUID NOT NULL UNIQUE,
    news_id UUID NOT NULL,
    PRIMARY KEY (comment_id,news_id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (news_id) REFERENCES news(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE news_comments
-- +goose StatementEnd
