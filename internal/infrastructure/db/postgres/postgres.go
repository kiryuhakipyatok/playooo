package db

import (
	"context"
	"crap/internal/config"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5"
)



func Connect(cfg *config.Config) (*pgx.Conn, error){
	url:=fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	conn,err:=pgx.Connect(ctx,url)
	if err!=nil{
		return nil,err
	}
	return conn,nil
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