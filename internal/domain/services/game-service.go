package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"slices"
)

type GameService interface {
	AddGameToUser(ctx context.Context, name, id string) error
	GetByName(ctx context.Context, name string) (*entities.Game, error)
	FetchGames(ctx context.Context, amount, page int) ([]entities.Game, error)
	DeleteGame(ctx context.Context, id, name string) error
}

type gameService struct {
	GameRepository repositories.GameRepository
	UserRepository repositories.UserRepository
	Transactor     repositories.Transactor
}

func NewGameService(gr repositories.GameRepository, ur repositories.UserRepository, t repositories.Transactor) GameService {
	return &gameService{
		GameRepository: gr,
		UserRepository: ur,
		Transactor:     t,
	}
}

func (gs *gameService) AddGameToUser(ctx context.Context, name, id string) error{
	_,err:=gs.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		game,err:=gs.GameRepository.FindByName(c,name)
		if err!=nil{
			return nil,err
		}
		user,err:=gs.UserRepository.FindById(c,id)
		if err!=nil{
			return nil,err
		}
		user.Games=append(user.Games, game.Name)
		game.NumberOfPlayers++
		if err:=gs.GameRepository.Save(c,*game);err!=nil{
			return nil,err
		}
		if err:=gs.UserRepository.Save(c,*user);err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (gs *gameService) GetByName(ctx context.Context, name string) (*entities.Game, error){
	game,err:=gs.GameRepository.FindByName(ctx,name)
	if err!=nil{
		return nil,err
	}
	return game,nil
}

func (gs *gameService) FetchGames(ctx context.Context, amount, page int) ([]entities.Game, error){
	games,err:=gs.GameRepository.Fetch(ctx,amount,page)
	if err!=nil{
		return nil,err
	}
	return games,nil
}

func (gs *gameService) DeleteGame(ctx context.Context, id, name string) error{
	_,err:=gs.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		game,err:=gs.GameRepository.FindByName(c,name)
		if err!=nil{
			return nil,err
		}
		user,err:=gs.UserRepository.FindById(c,id)
		if err!=nil{
			return nil,err
		}
		games:=slices.DeleteFunc(user.Games,func(g string) bool {
			return g == game.Name
		})
		user.Games = games
		game.NumberOfPlayers-- 
		if err:=gs.UserRepository.Save(c,*user);err!=nil{
			return nil,err
		}
		if err:=gs.GameRepository.Save(c,*game);err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}