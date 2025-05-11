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

type SheduleEvents struct {
	NoticeService services.NotificationService
	EventService  services.EventService
	UserService   services.UserService
	Logger        *logrus.Logger
	Bot           *bot.Bot
}

func (sheduleEvents *SheduleEvents) SetupSheduleEvents(stop chan struct{}) {
	cr := cron.New()
	cr.AddFunc("@every 1m", func() {
		now := time.Now()
		upcoming, err := sheduleEvents.EventService.FindUpcoming(context.Background(), now.Add(10*time.Minute).Add(15*time.Second))
		if err != nil {
			sheduleEvents.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range upcoming {
			if !event.NotifiedPre {
				premsg := "cобытие " + event.Body + " начнется через 10 минут!"
				sheduleEvents.NoticeService.CreateNotification(context.Background(), event, premsg)
				if err := sheduleEvents.Bot.SendMsg(event, premsg); err != nil {
					sheduleEvents.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
				}
				sheduleEvents.Logger.Infof("уведомление о предстоящем событии %v отправлено в %v", event.Body, time.Now())
				event.NotifiedPre = true
				if err := sheduleEvents.EventService.Save(context.Background(), event); err != nil {
					sheduleEvents.Logger.WithError(err).Errorf("failed to save event: %v", err)
				}
			}
		}
		current, err := sheduleEvents.EventService.FindUpcoming(context.Background(), now.Add(1*time.Minute).Add(15*time.Second))
		if err != nil {
			sheduleEvents.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range current {
			curmsg := "cобытие " + event.Body + " началось!"
			sheduleEvents.NoticeService.CreateNotification(context.Background(), event, curmsg)
			fmt.Printf("event: %s", event.Id)
			if err := sheduleEvents.Bot.SendMsg(event, curmsg); err != nil {
				sheduleEvents.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
			}
			sheduleEvents.Logger.Infof("уведомление о начале события %v отправлено в %v", event.Body, time.Now())
			if err := sheduleEvents.EventService.DeleteEvent(context.Background(), event.Id.String()); err != nil {
				sheduleEvents.Logger.WithError(err).Errorf("failed to delete event: %v", err)
			}
		}

	})
	cr.Start()
	<-stop
	if err := cr.Stop().Err(); err != nil {
		sheduleEvents.Logger.Infof("error to close sheduler: %v", err)
	} else {
		sheduleEvents.Logger.Info("stopping scheduler")
	}
}
