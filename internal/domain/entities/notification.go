package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id      uuid.UUID `json:"notice_id"`
	EventId uuid.UUID `json:"event_id"`
	Body    string    `json:"body"`
	Time time.Time `json:"time"`
}