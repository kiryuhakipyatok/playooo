package server

import (
	"context"
	"crap/internal/config"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg config.Config
	logger *logrus.Logger
}

func NewServer(cfg config.Config, logger *logrus.Logger) *Server{
	return &Server{
		cfg: cfg,
		logger: logger,
	}
}

func(s *Server) StartServer(stop chan struct{}){
	app:=fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))
	go func(){
		if err := app.Listen(s.cfg.Server.Host+s.cfg.Server.Port); err != nil {
			s.logger.Fatalf("failed to start server: %v", err)
		}
	}()
	s.logger.Info("server started succefully")
	<-stop
	shCtx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	if err:=app.ShutdownWithContext(shCtx);err!=nil{
		s.logger.Fatalf("server forced to shutdown: %v", err)
	}else{
		s.logger.Info("server stopped succefully")
	}
}