package sheduler

import (
	"context"
	"fmt"
	"crap/internal/sheduler/bot"
	"crap/internal/domain/services"
	"time"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Sheduler struct {
	NotificationService services.NotificationService
	EventService  services.EventService
	UserService   services.UserService
	Logger        *logrus.Logger
	Bot           *bot.Bot
}

func (s *Sheduler) SetupSheduler(stop chan struct{}) {
	if s.Bot != nil{
		s.Logger.Info("starting sheduller with bot")
	}else{
		s.Logger.Info("bot is nil, shedule whithout bot")
	}
	cr := cron.New()
	cr.AddFunc("@every 1m", func() {
		now := time.Now()
		upcoming, err := s.EventService.FindUpcoming(context.Background(), now.Add(10*time.Minute).Add(15*time.Second))
		if err != nil {
			s.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range upcoming {
			if !event.NotifiedPre {
				premsg := "cобытие " + event.Body + " начнется через 10 минут!"
				s.NotificationService.CreateNotification(context.Background(), event, premsg)
				if s.Bot!=nil{
					if err := s.Bot.SendMsg(event, premsg); err != nil {
						s.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
					}	
				}
				s.Logger.Infof("уведомление о предстоящем событии %v отправлено в %v", event.Body, time.Now())
				event.NotifiedPre = true
				if err := s.EventService.Save(context.Background(), event); err != nil {
					s.Logger.WithError(err).Errorf("failed to save event: %v", err)
				}
			}
		}
		current, err := s.EventService.FindUpcoming(context.Background(), now.Add(1*time.Minute).Add(15*time.Second))
		if err != nil {
			s.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range current {
			curmsg := "cобытие " + event.Body + " началось!"
			s.NotificationService.CreateNotification(context.Background(), event, curmsg)
			fmt.Printf("event: %s", event.Id)
			if s.Bot!=nil{
					if err := s.Bot.SendMsg(event, curmsg); err != nil {
				s.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
			}
			}
			s.Logger.Infof("уведомление о начале события %v отправлено в %v", event.Body, time.Now())
			if err := s.EventService.DeleteEvent(context.Background(), event.Id.String()); err != nil {
				s.Logger.WithError(err).Errorf("failed to delete event: %v", err)
			}
		}

	})
	cr.Start()
	<-stop
	if err := cr.Stop().Err(); err != nil {
		s.Logger.WithError(err).Info("error to close sheduler")
	} else {
		s.Logger.Info("stopping scheduler")
	}
}
