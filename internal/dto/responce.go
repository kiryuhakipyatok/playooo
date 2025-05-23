package dto

import (
	"github.com/google/uuid"
	"time"
)

type RegisterResponse struct {
	Id       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Telegram string    `json:"telegram"`
	Date time.Time `json:"date_of_register"`
}

type NewsResponse struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type EventResponse struct {
	Id       uuid.UUID `json:"id"`
	AuthorId uuid.UUID `json:"author-id"`
	Time     time.Time `json:"time"`
}
