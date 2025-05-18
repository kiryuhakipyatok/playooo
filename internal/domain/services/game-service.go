package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"crap/internal/dto"
	"slices"
)

type GameService interface {
	AddGameToUser(ctx context.Context, req dto.AddGameRequest) error
	GetByName(ctx context.Context, name string) (*entities.Game, error)
	FetchGames(ctx context.Context, req dto.PaginationRequest) ([]entities.Game, error)
	DeleteGame(ctx context.Context, req dto.DeleteGameRequest) error
	GetSorted(ctx context.Context, req dto.GamesSortRequest) ([]entities.Game, error)
	GetFiltered(ctx context.Context, req dto.GamesFilterRequest) ([]entities.Game, error)
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

func (gs *gameService) AddGameToUser(ctx context.Context, req dto.AddGameRequest) error{
	_,err:=gs.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		game,err:=gs.GameRepository.FindById(c,req.Game)
		if err!=nil{
			return nil,err
		}
		user,err:=gs.UserRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		user.Games=append(user.Games, game.Id)
		game.NumberOfPlayers++
		game.CalculateRating()
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

func (gs *gameService) FetchGames(ctx context.Context, req dto.PaginationRequest) ([]entities.Game, error){
	games,err:=gs.GameRepository.Fetch(ctx,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return games,nil
}

func (gs *gameService) DeleteGame(ctx context.Context, req dto.DeleteGameRequest) error{
	_,err:=gs.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		game,err:=gs.GameRepository.FindByName(c,req.Game)
		if err!=nil{
			return nil,err
		}
		user,err:=gs.UserRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		games:=slices.DeleteFunc(user.Games,func(g string) bool {
			return g == game.Name
		})
		user.Games = games
		game.NumberOfPlayers-- 
		game.CalculateRating()
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

func (gs *gameService) GetSorted(ctx context.Context, req dto.GamesSortRequest) ([]entities.Game, error){
	games,err:=gs.GameRepository.Sort(ctx,req.Field,req.Direction,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return games,nil
}

func (gs *gameService) GetFiltered(ctx context.Context, req dto.GamesFilterRequest) ([]entities.Game, error){
	games,err:=gs.GameRepository.Filter(ctx,req.Name,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return games,nil
}
