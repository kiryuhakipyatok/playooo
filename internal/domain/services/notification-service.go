package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, event entities.Event, msg string) error
	DeleteNotification(ctx context.Context, id string) error
	FetchNotifications(ctx context.Context, id string, amount, page int) ([]entities.Notification, error)
	DeleteAllNotifications(ctx context.Context, id string) error
}

type notificationService struct {
	NotificationRepository repositories.NotificationRepository
	EventRepository        repositories.EventRepository
	UserRepository         repositories.UserRepository
	Transactor             repositories.Transactor
}

func NewNotificationService(
	nr repositories.NotificationRepository,
	er repositories.EventRepository,
	ur repositories.UserRepository,
	t repositories.Transactor) NotificationService {
	return &notificationService{
		NotificationRepository: nr,
		EventRepository:  er,
		UserRepository:   ur,
		Transactor:       t,
	}
}

func (ns *notificationService) CreateNotification(ctx context.Context, event entities.Event, msg string) error{
	_,err:=ns.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		members,err:=ns.EventRepository.FetchMembers(c,event.Id.String())
		if err!=nil{
			return nil,err
		}
		for _,id:=range members{
			user,err:=ns.UserRepository.FindById(c,id)
			if err!=nil{
				return nil,err
			}
			notification:=entities.Notification{
				Id: uuid.New(),
				UserId: user.Id,
				EventId: event.Id,
				Body: msg,
				Time: time.Now(),
			}
			if err:=ns.NotificationRepository.Create(c,notification);err!=nil{
				return nil,err
			}
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (ns *notificationService) DeleteNotification(ctx context.Context, id string) error{
	_,err:=ns.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
		notification,err:=ns.NotificationRepository.FindById(c,id)
		if err!=nil{
			return nil,err
		}
		if err:=ns.NotificationRepository.DeleteById(c,notification.Id.String());err!=nil{
			return nil,err
		}
		return nil,nil
	})
	if err!=nil{
		return err
	}
	return nil
}

func (ns *notificationService) FetchNotifications(ctx context.Context, id string, amount, page int) ([]entities.Notification, error){
	user,err:=ns.UserRepository.FindById(ctx,id)
	if err!=nil{
		return nil,err
	}
	notifications,err:=ns.NotificationRepository.Fetch(ctx,user.Id.String(),amount,page)
	if err!=nil{
		return nil,err
	}
	return notifications,nil
}

func (ns *notificationService) DeleteAllNotifications(ctx context.Context, id string) error{
	user,err:=ns.UserRepository.FindById(ctx,id)
	if err!=nil{
		return err
	}
	if err:=ns.NotificationRepository.DeleteByUserId(ctx,user.Id.String());err!=nil{
		return err
	}
	return nil
}
