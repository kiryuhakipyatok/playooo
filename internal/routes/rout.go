package routes

import (
	"crap/internal/controllers/rest/handlers"

	"github.com/gofiber/fiber/v2"
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
	CommentsHandler    *handlers.CommentsHandler
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
	// rcfg.SetupSwaggerConfig()
}

func (rcfg *RoutConfig) SetupUserRoute() {
    userGroup := rcfg.App.Group("/api/users")

    userGroup.Get("/:id", rcfg.UserHandler.GetUser)
    userGroup.Get("", rcfg.UserHandler.GetUsers)

    userGroup.Patch("/avatar", rcfg.UserHandler.UploadAvatar)
    userGroup.Patch("/discord", rcfg.UserHandler.RecordDiscord)
    userGroup.Patch("/rating", rcfg.UserHandler.EditRating)

    userGroup.Delete("/avatar/:id", rcfg.UserHandler.DeleteAvatar)
}

func (rcfg *RoutConfig) SetupFriendshipsRoute() {
    friendshipGroup := rcfg.App.Group("/api/friends")

    friendshipGroup.Get("/requests", rcfg.FriendshipsHandler.GetFriendsRequests)
    friendshipGroup.Get("", rcfg.FriendshipsHandler.GetFriends)

    friendshipGroup.Patch("/accept", rcfg.FriendshipsHandler.AcceptFriendship)
    friendshipGroup.Patch("", rcfg.FriendshipsHandler.AddFriend)

    friendshipGroup.Delete("", rcfg.FriendshipsHandler.CancelFriendship)
}

func (rcfg *RoutConfig) SetupCommentRoute() {
    commentGroup := rcfg.App.Group("/api/comments")

    commentGroup.Get("", rcfg.CommentsHandler.GetComments)
    commentGroup.Post("", rcfg.CommentsHandler.AddComment)
}

func (rcfg *RoutConfig) SetupAuthRoute() {
    authGroup := rcfg.App.Group("/api/auth")

    authGroup.Get("/profile", rcfg.AuthHandler.Profile)

    authGroup.Post("/register", rcfg.AuthHandler.Register)
    authGroup.Post("/login", rcfg.AuthHandler.Login)
    authGroup.Post("/logout", rcfg.AuthHandler.Logout)
}

func (rcfg *RoutConfig) SetupGameRoute() {
    gameGroup := rcfg.App.Group("/api/games")

    gameGroup.Get("/sort", rcfg.GameHandler.GetSortedGames)
    gameGroup.Get("/filter", rcfg.GameHandler.GetFilteredGames)
    gameGroup.Get("/:name", rcfg.GameHandler.GetGame)
    gameGroup.Get("", rcfg.GameHandler.GetGames)

    gameGroup.Patch("", rcfg.GameHandler.AddGame)

    gameGroup.Delete("", rcfg.GameHandler.DeleteGame)
}

func (rcfg *RoutConfig) SetupEventRoute() {
    eventsGroup := rcfg.App.Group("/api/events")

    eventsGroup.Get("/sort", rcfg.EventHandler.GetSortedEvents)
    eventsGroup.Get("/filter", rcfg.EventHandler.GetFilteredEvents)
    eventsGroup.Get("/:id", rcfg.EventHandler.GetEvent)
    eventsGroup.Get("", rcfg.EventHandler.GetEvents)

    eventsGroup.Patch("/join", rcfg.EventHandler.Join)
    eventsGroup.Patch("/unjoin", rcfg.EventHandler.Unjoin)

    eventsGroup.Post("", rcfg.EventHandler.CreateEvent)
}

func (rcfg *RoutConfig) SetupNewsRoute() {
    newsGroup := rcfg.App.Group("/api/news")

    newsGroup.Get("/:id", rcfg.NewsHandler.GetNews)
    newsGroup.Get("", rcfg.NewsHandler.GetSomeNews) // Fixed duplicate "/api/news" prefix

    newsGroup.Post("", rcfg.NewsHandler.CreateNews)
}

func (cfg *RoutConfig) SetupNotificationsRoute() {
    notificationsGroup := cfg.App.Group("/api/notifications")

    notificationsGroup.Get("", cfg.NoticeHandler.GetNotifications)

    notificationsGroup.Delete("/all/:id", cfg.NoticeHandler.DeleteAllNotifications) 
    notificationsGroup.Delete("", cfg.NoticeHandler.DeleteNotification)
}

// func (cfg *RoutConfig) SetupSwaggerConfig() {
// 	cfg.App.Get("/swagger/*", swagger.HandlerDefault)
// }
