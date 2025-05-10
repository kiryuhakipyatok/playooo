package entities

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id          uuid.UUID      `json:"event_id"`
	AuthorId    uuid.UUID      `json:"author_id"`
	Body        string         `json:"body"`
	Game        string         `json:"game"`
	Max         int            `json:"max"`
	Time        time.Time      `json:"minute"`
	NotifiedPre bool
}