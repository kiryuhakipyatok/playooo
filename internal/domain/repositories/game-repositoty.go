package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"crap/internal/domain/entities"
)

type GameRepository interface {
	Save(ctx context.Context, game entities.Game) error
	FindByName(ctx context.Context, name string) (*entities.Game, error)
	Fetch(ctx context.Context, amount, page int) ([]entities.Game, error)
}

type gameRepository struct {
	DB *pgx.Conn
}

func NewGameRepository(db *pgx.Conn) GameRepository {
	return &gameRepository{
		DB: db,
	}
}

func (gr *gameRepository) Save(ctx context.Context, game entities.Game) error{
	if _,err:=gr.DB.Exec(ctx,"UPDATE games SET name=$1,description=$2,banner=$3,picture=$4,number_of_players=$5,number_of_events=$6,rating=$7 WHERE id=$8",game.Name,game.Description,game.Banner,game.Picture,game.NumberOfPlayers,game.NumberOfEvents,game.Rating,game.Banner,game.Picture,game.Description,game.Id);err!=nil{
		return err
	}
	return nil
}

func (gr *gameRepository) FindByName(ctx context.Context, id string) (*entities.Game, error){
	game:=entities.Game{}
	if err:=gr.DB.QueryRow(ctx,"SELECT * FROM games WHERE id = $1",id).Scan(&game.Id,&game.Name,&game.Description,&game.Banner,&game.Picture,&game.NumberOfPlayers,&game.NumberOfEvents,&game.Rating);err!=nil{
		return nil,err
	}
	return &game, nil
}

func (gr *gameRepository) Fetch(ctx context.Context, amount, page int) ([]entities.Game, error){
	games := []entities.Game{}
	rows,err:=gr.DB.Query(ctx,"SELECT * FROM games ORDER BY number_of_players DESC OFFSET $1 LIMIT $2",amount*page-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		game:=entities.Game{}
		if err:=rows.Scan(&game.Id,&game.Name,&game.Description,&game.Banner,&game.Picture,&game.NumberOfPlayers,&game.NumberOfEvents,&game.Rating);err!=nil{
			return nil,err
		}
		games=append(games, game)
	}
	return games, nil
}
