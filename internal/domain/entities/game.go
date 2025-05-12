package entities

import "github.com/google/uuid"

type Game struct {
	Id              uuid.UUID `json:"id"`
	Name            string  `json:"name"`
	NumberOfPlayers int     `json:"num_of_players"`
	NumberOfEvents  int     `json:"num_of_events"`
	Rating          float64 `json:"rating"`
}

func (g *Game) CalculateRating() {
	if g.NumberOfPlayers == 0 {
		g.Rating = 0
		return
	}
	g.Rating = float64(g.NumberOfPlayers+g.NumberOfEvents) / 2.0
}