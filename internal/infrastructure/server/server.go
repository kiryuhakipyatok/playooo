package server

import (
	"context"
	"crap/internal/config"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/gofiber/jwt/v3"
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
	app.Static("files", "../../files")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}),func (c *fiber.Ctx) error {
		if c.Path() == "/api/register" || c.Path() == "/api/login"{
			return c.Next()
		}
		return jwtware.New(jwtware.Config{
			SigningKey: []byte(s.cfg.Auth.Secret),
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(fiber.Map{
					"message": "unauthorized",
				})
			},
			TokenLookup: "cookie:jwt",
		})(c)
	})
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