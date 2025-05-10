package entities

import (
	"github.com/google/uuid"
	"time"
)

type News struct {
	Id       uuid.UUID      `json:"news_id"`
	Title    string         `json:"title"`
	Body     string         `json:"body"`
	Time     time.Time      `json:"time"`
	Link     string         `json:"link"`
	Picture  string         `json:"picture"`
}
