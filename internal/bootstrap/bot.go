package bootstrap

import (
	"crap/internal/config"
	"crap/internal/domain/repositories"
	"crap/internal/sheduler/bot"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func StartBot(cfg config.Config,db *pgx.Conn, redis *redis.Client,l *logrus.Logger, stop chan struct{}) (*bot.Bot, error) {
	userRepository := repositories.NewUserRepository(db, redis)
	eventRepository := repositories.NewEventRepository(db, redis)
	bot, err := bot.CreateBot(stop,l, userRepository, eventRepository,cfg.Bot.Token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
