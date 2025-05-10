package entities

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Id           uuid.UUID `json:"comment_id"`
	AuthorId     uuid.UUID `json:"author_id"`
	Body         string    `json:"body"`
	Time         time.Time `json:"time"`
}