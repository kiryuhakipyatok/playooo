package server

import (
	"crap/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/gofiber/contrib/websocket"
	"crap/internal/chat"
)

func CreateServer(cfg *config.Config) (*fiber.App, error) {
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/ws/:id", websocket.New(chat.HandleWS))
	app.Static("files", "../../files")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}), func(c *fiber.Ctx) error {
		excludedPaths := map[string]bool{
			"/api/auth/register": true,
			"/api/auth/login":    true,
		}

		if excludedPaths[c.Path()] {
			return c.Next()
		}

		return jwtware.New(jwtware.Config{
			SigningKey: []byte(cfg.Auth.Secret),
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(fiber.Map{
					"message": "unauthorized",
				})
			},
			TokenLookup: "cookie:jwt",
		})(c)
	})
	app.Use("/ws/:id", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	return app, nil
}
