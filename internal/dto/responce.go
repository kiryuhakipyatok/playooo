package dto

import (
	"github.com/google/uuid"
	"time"
)

type RegisterResponse struct {
	Id       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Telegram string    `json:"telegram"`
}

type NewsResponse struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type EventResponse struct {
	Id       uuid.UUID `json:"id"`
	AuthorId uuid.UUID `json:"author-id"`
	Body     string    `json:"body"`
	Game     string    `json:"game"`
	Max      int       `json:"max"`
	Time     time.Time `json:"time"`
}
