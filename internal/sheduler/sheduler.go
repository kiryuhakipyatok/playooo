package sheduler

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/sheduler/bot"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Sheduler struct {
	NotificationService services.NotificationService
	EventService        services.EventService
	UserService         services.UserService
	Logger              *logrus.Logger
	Bot                 *bot.Bot
}

func (s *Sheduler) SetupSheduler(stop chan struct{}) {
	if s.Bot != nil {
		s.Logger.Info("starting sheduller with bot")
	} else {
		s.Logger.Info("bot is nil, shedule whithout bot")
	}
	cr := cron.New()
	if _, err := cr.AddFunc("@every 1m", func() {
		now := time.Now()
		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		upcoming, err := s.EventService.FindUpcoming(ctx1, now.Add(10*time.Minute).Add(30*time.Second))
		if err != nil {
			s.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range upcoming {
			if !event.NotificatedPre {
				premsg := "cобытие " + event.Body + " начнется через 10 минут!"
				if err := s.NotificationService.CreateNotification(ctx1, event, premsg); err != nil {
					s.Logger.WithError(err).Errorf("failed to create notification: %v", err)
				}
				if s.Bot != nil {
					if err := s.Bot.SendMsg(event, premsg); err != nil {
						s.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
					}
				}
				s.Logger.Infof("уведомление о предстоящем событии %v отправлено в %v", event.Body, time.Now())
				event.NotificatedPre = true
				if err := s.EventService.Save(context.Background(), event); err != nil {
					s.Logger.WithError(err).Errorf("failed to save event: %v", err)
				}
			}
		}
		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		current, err := s.EventService.FindUpcoming(ctx2, now.Add(30*time.Second))
		if err != nil {
			s.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range current {
			curmsg := "cобытие " + event.Body + " началось!"
			if err := s.NotificationService.CreateNotification(ctx2, event, curmsg); err != nil {
				s.Logger.WithError(err).Errorf("failed to create notification: %v", err)
			}
			if s.Bot != nil {
				if err := s.Bot.SendMsg(event, curmsg); err != nil {
					s.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
				}
			}
			s.Logger.Infof("уведомление о начале события %v отправлено в %v", event.Body, time.Now())
			if err := s.EventService.DeleteEvent(context.Background(), event.Id.String()); err != nil {
				s.Logger.WithError(err).Errorf("failed to delete event: %v", err)
			}
		}
	}); err != nil {
		s.Logger.WithError(err).Error("failed to add cron job")
		return
	}
	cr.Start()
	<-stop
	if err := cr.Stop().Err(); err != nil {
		s.Logger.WithError(err).Error("error stopping scheduler")
	} else {
		s.Logger.Info("scheduler stopped successfully")
	}
}
