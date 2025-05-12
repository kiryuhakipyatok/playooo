package bootstrap

import (
	"crap/internal/config"
	"crap/internal/controllers/rest/handlers"
	"crap/internal/domain/repositories"
	"crap/internal/domain/services"
	"crap/internal/routes"
	"crap/internal/sheduler"
	"crap/internal/sheduler/bot"

	// "crap/internal/sheduler"
	// "crap/internal/sheduler/bot"
	// "sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	App *fiber.App
	Postgres  *pgx.Conn
	Redis     *redis.Client
	Logger    *logrus.Logger
	Validator *validator.Validate
}

func NewBootstrapConfig(a *fiber.App,p *pgx.Conn, r *redis.Client, l *logrus.Logger, v *validator.Validate) BootstrapConfig{
	return BootstrapConfig{
		App:a,
		Postgres: p,
		Redis: r,
		Logger: l,
		Validator: v,
	}
}

func(bcfg *BootstrapConfig) BootstrapHandlers(stop chan struct{}, cfg *config.Config) {

	transactor := repositories.NewTransactor(bcfg.Postgres)

	userRepository := repositories.NewUserRepository(bcfg.Postgres, bcfg.Redis)
	gameRepository := repositories.NewGameRepository(bcfg.Postgres)
	eventRepository := repositories.NewEventRepository(bcfg.Postgres, bcfg.Redis)
	newsRepository := repositories.NewNewsRepository(bcfg.Postgres)
	commentRepository := repositories.NewCommentRepository(bcfg.Postgres)
	notificationRepository := repositories.NewNoticeRepository(bcfg.Postgres)
	friendshipsRepository:=repositories.NewFriendshipsRepository(bcfg.Postgres)

	userService := services.NewUserService(userRepository, transactor,cfg)
	authService := services.NewAuthService(userRepository,cfg)
	gameService := services.NewGameService(gameRepository, userRepository, transactor)
	eventService := services.NewEventService(eventRepository, userRepository, gameRepository, transactor)
	newsService := services.NewNewsService(newsRepository, transactor,cfg)
	notificationService := services.NewNotificationService(notificationRepository, eventRepository, userRepository, transactor)
	commentService := services.NewCommentService(commentRepository, userRepository, eventRepository, newsRepository, transactor)
	friendshipsService :=services.NewFriendshipsService(friendshipsRepository,userRepository)

	userHandler := handlers.NewUsersHandler(userService, bcfg.Logger, bcfg.Validator)
	authHander := handlers.NewAuthHandler(authService, bcfg.Logger, bcfg.Validator)
	gameHandler := handlers.NewGamesHandler(gameService, bcfg.Logger, bcfg.Validator)
	eventHandler := handlers.NewEventHandler(eventService, bcfg.Logger, bcfg.Validator)
	newsHandler := handlers.NewNewsHandler(newsService, bcfg.Logger, bcfg.Validator)
	notificationHandler := handlers.NewNotificationsHandler(notificationService, bcfg.Logger, bcfg.Validator)
	commetHandler :=handlers.NewCommentHandler(commentService,bcfg.Logger,bcfg.Validator)
	friendshipsHandler:=handlers.NewFriendshipsHandler(friendshipsService,bcfg.Logger,bcfg.Validator)

	routConfig := routes.RoutConfig{
		App:           bcfg.App,
		UserHandler:   &userHandler,
		AuthHandler:   &authHander,
		GameHandler:   &gameHandler,
		EventHandler:  &eventHandler,
		NewsHandler:   &newsHandler,
		NoticeHandler: &notificationHandler,
		CommentsHandler: &commetHandler,
		FriendshipsHandler: &friendshipsHandler,
	}

	routConfig.Setup()
}

func(bcfg *BootstrapConfig) BootstrapSheduler(stop chan struct{}, bot *bot.Bot, cfg *config.Config) sheduler.Sheduler{
	transactor := repositories.NewTransactor(bcfg.Postgres)
	userRepository := repositories.NewUserRepository(bcfg.Postgres, bcfg.Redis)
	userService := services.NewUserService(userRepository, transactor,cfg)
	eventRepository := repositories.NewEventRepository(bcfg.Postgres, bcfg.Redis)
	gameRepository := repositories.NewGameRepository(bcfg.Postgres)
	eventService := services.NewEventService(eventRepository, userRepository, gameRepository, transactor)
	notificationRepository := repositories.NewNoticeRepository(bcfg.Postgres)
	notificationService := services.NewNotificationService(notificationRepository, eventRepository, userRepository, transactor)
	sheduler:=sheduler.Sheduler{
		NotificationService: notificationService,
		UserService: userService,
		EventService: eventService,
		Logger: bcfg.Logger,
		Bot: bot,
	}
	return sheduler
}

func(bcfg *BootstrapConfig) BootstrapBot(stop chan struct{}, cfg *config.Config) (*bot.Bot,error){
	userRepository := repositories.NewUserRepository(bcfg.Postgres, bcfg.Redis)
	eventRepository := repositories.NewEventRepository(bcfg.Postgres, bcfg.Redis)
	bot, err := bot.CreateBot(stop,bcfg.Logger, userRepository, eventRepository,cfg.Bot.Token)
	if err != nil {
		return nil,err
	}
	return bot, nil
}
