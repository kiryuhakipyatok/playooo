package db

import (
	"context"
	"crap/internal/config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)





func Connect(cfg config.Config) (*redis.Client,error){
	addr:=fmt.Sprintf("%s:%s",cfg.Redis.Host,cfg.Redis.Port)
	client:=redis.NewClient(&redis.Options{
		Addr: addr,
		Password: cfg.Redis.Password,
		DB: 0,
	})
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	if _,err:=client.Ping(ctx).Result();err!=nil{
		return nil,err
	}
	return client,nil
}