package server

import (
	"crap/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
)

func CreateServer(cfg config.Config) (*fiber.App,error){
	app:=fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Static("files", "../../files")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}),func (c *fiber.Ctx) error {
		excludedPaths := map[string]bool{
            "/api/register": true,
            "/api/login":    true,
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
	return app,nil
}