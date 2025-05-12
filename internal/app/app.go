package app

import (
	"context"
	"fmt"
	"sync"

	// "crap/internal/bootstrap"
	"crap/internal/bootstrap"
	"crap/internal/config"
	p "crap/internal/infrastructure/db/postgres"
	r "crap/internal/infrastructure/db/redis"
	"crap/internal/infrastructure/server"
	"crap/pkg/logger"
	"crap/pkg/validator"
	"os"
	"os/signal"
	"syscall"

	// "crap/internal/sheduler/bot"
	"time"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	logger := logger.NewLogger()
	validator := validator.NewValidator()
	stop := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	postgres, err := p.Connect(cfg)
	if err != nil {
		logger.Fatal("failed to connect to postgres")
	} else {
		logger.Info("connect to postgres succefully")
	}
	redis, err := r.Connect(cfg)
	if err != nil {
		logger.Info("failed to connect to redis")
	} else {
		logger.Info("connect to redis succefully")
	}
	app, err := server.CreateServer(cfg)
	if err != nil {
		logger.WithError(err).Fatal("failed to create server")
	} else {
		logger.Info("server created succefully")
	}
	bcfg := bootstrap.NewBootstrapConfig(app, postgres, redis, logger, validator)
	bcfg.BootstrapHandlers(stop, cfg)
	bot,err:=bcfg.BootstrapBot(stop,cfg)
	sheduler:=bcfg.BootstrapSheduler(stop,bot)
	if err!=nil{
		logger.WithError(err).Info("error start bot")
	}else{
		logger.Info("bot started successful")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Listen("0.0.0.0"+ ":" + cfg.Server.Port); err != nil {
			logger.WithError(err).Fatal("failed to start server")
		}
	}()
	wg.Add(1)
	go func(){
		defer wg.Done()
		bot.ListenForUpdates(stop)
	}()
	wg.Add(1)
	go func(){
		defer wg.Done()
		sheduler.SetupSheduler(stop)
	}()
	<-quit
	close(quit)
	close(stop)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.WithError(err).Fatal("server forced to shutdown")
	}
	if err := postgres.Close(ctx); err != nil {
		logger.WithError(err).Fatal("failed to close postgres")
	} else {
		logger.Info("close postgres success")
	}
	if redis != nil {
		if err := redis.Close(); err != nil {
			logger.WithError(err).Fatal("failed to close redis")
		} else {
			logger.Info("close redis success")
		}
	}
	logger.Info("server stopped successfuly")
}
