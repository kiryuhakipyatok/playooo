package app

import (
	"crap/internal/config"
	pg "crap/internal/infrastructure/db/postgres"
	r "crap/internal/infrastructure/db/redis"
	"crap/internal/infrastructure/server"
	"crap/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	// "context"
	// "time"
)

func Run() {
	cfg, err := config.LoadConfig("config")
	logger:=logger.NewLogger()
	if err!=nil{
		panic(err)
	}
	server:=server.NewServer(cfg,logger)
	stop:=make(chan struct{},1)
	quit:=make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(3)
	go func(){
		defer wg.Done()
		r.NewRedis(cfg,logger).Connect(stop)
	}()
	go func(){
		defer wg.Done()
		pg.NewPostgres(cfg,logger).Connect(stop)
	}()
	go func(){
		defer wg.Done()
		server.StartServer(stop)
	}()
	<-quit
	close(stop)
	wg.Wait()
}