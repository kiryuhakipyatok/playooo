package entities

import (
	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID      `json:"id"`
	Login           string         `json:"login"`
	Telegram      	string         `json:"telegram"`
	ChatId          string         `json:"chat_id"`
	Rating          float64        `json:"rating"`
	TotalRating     int            `json:"total_rating"`
	NumberOfRatings int            `json:"number_of_ratings"`
	Games           []string 	`json:"games"`
	Password        []byte	       `json:"-"`
	Avatar          string		`json:"avatar"`
	Discord         string		`json:"discord"`
}