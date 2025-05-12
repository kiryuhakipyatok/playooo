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
	if _,err:=gr.DB.Exec(ctx,"UPDATE games SET name=$1,number_of_players=$2,number_of_events=$3,rating=$4 WHERE name=$5",game.Name,game.NumberOfPlayers,game.NumberOfEvents,game.Rating,game.Name);err!=nil{
		return err
	}
	return nil
}

func (gr *gameRepository) FindByName(ctx context.Context, name string) (*entities.Game, error){
	game:=entities.Game{}
	if err:=gr.DB.QueryRow(ctx,"SELECT * FROM games WHERE name = $1",name).Scan(&game.Id,&game.Name,&game.NumberOfPlayers,&game.NumberOfEvents,&game.Rating);err!=nil{
		return nil,err
	}
	return &game, nil
}

func (gr *gameRepository) Fetch(ctx context.Context, amount, page int) ([]entities.Game, error){
	games := []entities.Game{}
	rows,err:=gr.DB.Query(ctx,"SELECT * FROM games OFFSET $1 LIMIT $2",amount*page-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		game:=entities.Game{}
		if err:=rows.Scan(&game.Name,&game.NumberOfPlayers,&game.NumberOfEvents,&game.Rating);err!=nil{
			return nil,err
		}
		games=append(games, game)
	}
	return games, nil
}
