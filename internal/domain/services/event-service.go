package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"crap/internal/dto"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context,req dto.CreateEventRequest) (*entities.Event, error)
	GetById(ctx context.Context, id string) (*entities.Event, error)
	FetchEvents(ctx context.Context, req dto.PaginationRequest) ([]entities.Event, error)
	FindUpcoming(ctx context.Context, time time.Time) ([]entities.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	Save(ctx context.Context, event entities.Event) error
	Join(ctx context.Context, req dto.JoinToEventRequest) error
	Unjoin(ctx context.Context, req dto.UnjoinFromEventRequest) error
	GetSorted(ctx context.Context, req dto.EventsSortRequest) ([]entities.Event, error)
	GetFiltered(ctx context.Context, req dto.EventsFilterRequest) ([]entities.Event, error)
}

type eventService struct {
	EventRepository repositories.EventRepository
	UserRepository  repositories.UserRepository
	GameRepository  repositories.GameRepository
	Transactor      repositories.Transactor
}

func NewEventService(
	eventRepository repositories.EventRepository,
	userRepository repositories.UserRepository,
	gameRepository repositories.GameRepository,
	transactor repositories.Transactor) EventService {
	return &eventService{
		EventRepository: eventRepository,
		UserRepository:  userRepository,
		GameRepository:  gameRepository,
		Transactor:      transactor,
	}
}

func (es *eventService)	CreateEvent(ctx context.Context, req dto.CreateEventRequest) (*entities.Event, error){
	res,err:=es.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		user,err:=es.UserRepository.FindById(c,req.AuthorId)
		if err!=nil{
			return nil,err
		}
		game,err:=es.GameRepository.FindById(c,req.Game)
		if err!=nil{
			return nil,err
		}
		if !slices.Contains(user.Games,game.Id){
			return nil,errors.New("user does not have this game")
		}
		event:=entities.Event{
			Id: uuid.New(),
			AuthorId: user.Id,
			Body: req.Body,
			Game: game.Name,
			Max: req.Max,
			Time: time.Now().Add(time.Minute*time.Duration(req.Minute)),
		}
		if req.Minute <= 10 {
			event.NotificatedPre = true
		}
		if err:=es.EventRepository.Create(c,event);err!=nil{
			return nil,err
		}
		game.NumberOfEvents++
		game.CalculateRating()
		if err:=es.GameRepository.Save(c,*game);err!=nil{
			return nil,err
		}
		return event,nil
	})
	if err!=nil{
		return nil,err
	}
	return res.(*entities.Event),nil
}

func (es *eventService)	GetById(ctx context.Context, id string) (*entities.Event, error){
	event,err:=es.EventRepository.FindById(ctx,id)
	if err!=nil{
		return nil,err
	}
	return event,nil
}

func (es *eventService)	FetchEvents(ctx context.Context, req dto.PaginationRequest) ([]entities.Event, error){
	events,err:=es.EventRepository.Fetch(ctx,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return events,nil
}

func (es *eventService)	FindUpcoming(ctx context.Context, time time.Time) ([]entities.Event, error){
	events,err:=es.EventRepository.FetchUpcoming(ctx,time)
	if err!=nil{
		return nil,err
	}
	return events,nil
}

func (es eventService) Save(c context.Context, event entities.Event) error {
	if err := es.EventRepository.Save(c, event); err != nil {
		return err
	}
	return nil
}

func (es *eventService)	DeleteEvent(ctx context.Context, id string) error{
	_,err:=es.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		event,err:=es.EventRepository.FindById(c,id)
		if err!=nil{
			return nil,err
		}
		if err:=es.EventRepository.Delete(c,*event);err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (es *eventService)	Join(ctx context.Context, req dto.JoinToEventRequest) error{
	_,err:=es.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		user,err:=es.UserRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		event,err:=es.EventRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		if err:=es.EventRepository.Join(c,user.Id.String(),event.Id.String());err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (es *eventService)	Unjoin(ctx context.Context, req dto.UnjoinFromEventRequest) error{
	_,err:=es.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		user,err:=es.UserRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		event,err:=es.EventRepository.FindById(c,req.EventId)
		if err!=nil{
			return nil,err
		}
		if err:=es.EventRepository.Unjoin(c,user.Id.String(),event.Id.String());err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (es *eventService) GetSorted(ctx context.Context, req dto.EventsSortRequest) ([]entities.Event, error){
	events,err:=es.EventRepository.Sort(ctx,req.Field,req.Direction,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return events,nil
}

func (es *eventService)	GetFiltered(ctx context.Context, req dto.EventsFilterRequest) ([]entities.Event, error){
	events,err:=es.EventRepository.Filter(ctx,req.Game,req.Max,req.Time,req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	return events,nil
}
