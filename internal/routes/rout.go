package routes

import (
	"crap/internal/controllers/rest/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RoutConfig struct {
	App                *fiber.App
	UserHandler        *handlers.UsersHandler
	AuthHandler        *handlers.AuthHandler
	GameHandler        *handlers.GamesHandler
	EventHandler       *handlers.EventsHandler
	NewsHandler        *handlers.NewsHandler
	NoticeHandler      *handlers.NotificationsHandler
	FriendshipsHandler *handlers.FriendshipsHandler
	CommentsHandler *handlers.CommentsHandler
}

func (rcfg *RoutConfig) Setup() {
	rcfg.SetupUserRoute()
	rcfg.SetupAuthRoute()
	rcfg.SetupGameRoute()
	rcfg.SetupNewsRoute()
	rcfg.SetupEventRoute()
	rcfg.SetupNotificationsRoute()
	rcfg.SetupFriendshipsRoute()
	rcfg.SetupCommentRoute()
	rcfg.SetupSwaggerConfig()
}

func (rcfg *RoutConfig) SetupUserRoute() {
	userGroup:=rcfg.App.Group("/api/users")

	userGroup.Patch("/avatar/:id", rcfg.UserHandler.UploadAvatar)
	userGroup.Patch("/discord", rcfg.UserHandler.RecordDiscord)
	userGroup.Patch("/rating", rcfg.UserHandler.EditRating)

	userGroup.Delete("/avatar/:id", rcfg.UserHandler.DeleteAvatar)

	userGroup.Get("/:id", rcfg.UserHandler.GetUser)
	userGroup.Get("", rcfg.UserHandler.GetUsers)
}

func (rcfg *RoutConfig) SetupFriendshipsRoute(){
	friendshipGroup:=rcfg.App.Group("/api/friends")

	friendshipGroup.Patch("", rcfg.FriendshipsHandler.AddFriend)
	friendshipGroup.Patch("/accept", rcfg.FriendshipsHandler.AcceptFriendship)

	friendshipGroup.Delete("", rcfg.FriendshipsHandler.CancelFriendship)

	friendshipGroup.Get("",rcfg.FriendshipsHandler.GetFriends)
	friendshipGroup.Get("/requests",rcfg.FriendshipsHandler.GetFriendsRequests)
}

func(rcfg *RoutConfig) SetupCommentRoute(){
	commentGroup:=rcfg.App.Group("/api/comments")

	commentGroup.Post("", rcfg.CommentsHandler.AddComment)

	commentGroup.Get("", rcfg.CommentsHandler.GetComments)
}

func (rcfg *RoutConfig) SetupAuthRoute() {
	rcfg.App.Post("/api/register", rcfg.AuthHandler.Register)
	rcfg.App.Post("/api/login", rcfg.AuthHandler.Login)
	rcfg.App.Post("/api/logout", rcfg.AuthHandler.Logout)

	rcfg.App.Get("/api/profile", rcfg.AuthHandler.Profile)
}

func (rcfg *RoutConfig) SetupGameRoute() {
	gameGroup := rcfg.App.Group("/api/games")

	gameGroup.Patch("", rcfg.GameHandler.AddGame)

	gameGroup.Delete("", rcfg.GameHandler.DeleteGame)

	gameGroup.Get("/:name", rcfg.GameHandler.GetGame)
	gameGroup.Get("", rcfg.GameHandler.GetGames)
}

func (rcfg *RoutConfig) SetupEventRoute() {
	eventsGroup := rcfg.App.Group("/api/events")
	eventsGroup.Post("", rcfg.EventHandler.CreateEvent)

	eventsGroup.Patch("/join", rcfg.EventHandler.Join)
	eventsGroup.Patch("/unjoin", rcfg.EventHandler.Unjoin)

	eventsGroup.Get("/:id", rcfg.EventHandler.GetEvent)
	eventsGroup.Get("/api/events", rcfg.EventHandler.GetEvents)
}

func (rcfg *RoutConfig) SetupNewsRoute() {
	newsGroup := rcfg.App.Group("/api/news")

	newsGroup.Post("", rcfg.NewsHandler.CreateNews)

	newsGroup.Get("/:id", rcfg.NewsHandler.GetNews)
	newsGroup.Get("/api/news", rcfg.NewsHandler.GetSomeNews)
}

func (cfg *RoutConfig) SetupNotificationsRoute() {
	notificationsGroup := cfg.App.Group("/api/notifications")

	notificationsGroup.Delete("/:id", cfg.NoticeHandler.DeleteNotification)
	notificationsGroup.Delete("/:id", cfg.NoticeHandler.DeleteAllNotifications)

	notificationsGroup.Get("", cfg.NoticeHandler.GetNotifications)
}

func (cfg *RoutConfig) SetupSwaggerConfig() {
	cfg.App.Get("/swagger/*", swagger.HandlerDefault)
}
