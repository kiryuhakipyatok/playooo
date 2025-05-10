package db

import (
	"context"
	"crap/internal/config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	cfg config.Config
	logger *logrus.Logger
}

func NewRedis(cfg config.Config, logger *logrus.Logger) *Redis{
	return &Redis{
		cfg: cfg,
		logger: logger,
	}
}

func(r *Redis) Connect(stop chan struct{}){
	addr:=fmt.Sprintf("%s:%s",r.cfg.Redis.Host,r.cfg.Redis.Port)
	client:=redis.NewClient(&redis.Options{
		Addr: addr,
		Password: r.cfg.Redis.Password,
		DB: 0,
	})
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	if _,err:=client.Ping(ctx).Result();err!=nil{
		r.logger.Info("error to connect to redis")
	}else{
		r.logger.Info("connect to redis successfully")
		<-stop
		if err:=client.Close();err!=nil{
			r.logger.Info("error to close redis")
		}else{
			r.logger.Info("redis closed successfully")
		}
	}
	
}