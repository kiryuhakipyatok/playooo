package app

import (
	"crap/internal/config"
	db "crap/internal/infrastructure/db/postgres"
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
	wg.Add(2)
	go func(){
		defer wg.Done()
		db.NewPostgres(cfg,logger).Connect(stop)
	}()
	go func(){
		defer wg.Done()
		server.StartServer(stop)
	}()
	<-quit
	close(stop)
	wg.Wait()
}