package db

import (
	"context"
	"crap/internal/config"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	cfg config.Config
	logger *logrus.Logger
}

func NewPostgres(cfg config.Config,logger *logrus.Logger) *Postgres{
	return &Postgres{
		cfg: cfg,
		logger: logger,
	}
}

func(pg *Postgres) Connect(stop chan struct{}){
	url:=fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.cfg.Postgres.User,
		pg.cfg.Postgres.Password,
		pg.cfg.Postgres.Host,
		pg.cfg.Postgres.Port,
		pg.cfg.Postgres.Database,
	)
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	conn,err:=pgx.Connect(ctx,url)
	if err!=nil{
		pg.logger.Fatalf("error to connect to postgres: %v",err)
	}
	pg.logger.Info("connect to postgres succesfully")
	<-stop
	if err:=conn.Close(ctx);err!=nil{
		pg.logger.Fatalf("error to close postgres: %v",err)
	}
	pg.logger.Info("close postgres succesfully")
}

// func(pg *Postgres) CloseConn(stop chan struct{},conn *pgx.Conn){
// 	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
// 	defer cancel()
// 	<-stop
// 	if err:=conn.Close(ctx);err!=nil{
// 		pg.logger.Fatalf("error to close postgres connection: %v",err)
// 	}
// 	pg.logger.Info("close postgres succesfully")
// }