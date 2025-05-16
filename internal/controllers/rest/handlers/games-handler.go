package handlers

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type GamesHandler struct {
	GameService services.GameService
	Logger      *logrus.Logger
	Validator   *validator.Validate
}

func NewGamesHandler(fs services.GameService, l *logrus.Logger, v *validator.Validate) GamesHandler {
	return GamesHandler{
		GameService: fs,
		Logger:      l,
		Validator:   v,
	}
}

// AddGame godoc
// @Summary Add a game to user
// @Description Add a game to user's collection
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AddGameRequest true "Add Game Request"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games [post]
func(gh *GamesHandler) AddGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "add-game")
	request := dto.AddGameRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := gh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err := gh.GameService.AddGameToUser(ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to add game: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// DeleteGame godoc
// @Summary Delete a game from user
// @Description Remove a game from user's collection
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.DeleteGameRequest true "Delete Game Request"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games [delete]
func(gh *GamesHandler) DeleteGame(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "delete-game")
	request := dto.DeleteGameRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := gh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=gh.GameService.DeleteGame(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete game: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// GetGames godoc
// @Summary Get games
// @Description Retrieve paginated list of games
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.PaginationRequest true "Pagination parameters"
// @Success 200 {array} entities.Game "List of games"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games [get]
func(gh *GamesHandler) GetGames(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-games")
	params := dto.PaginationRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=gh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	games,err:=gh.GameService.FetchGames(ctx,params)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get games: " + err.Error(),
		})
	}
	return c.JSON(games)
}

// GetGame godoc
// @Summary Get game details
// @Description Get detailed information about specific game
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param game path string true "Game title"
// @Success 200 {object} entities.Game "Game data"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games/{game} [get]
func(gh *GamesHandler) GetGame(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-game")
	title:=c.Params("game")
	game,err:=gh.GameService.GetByName(ctx,title)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get game: " + err.Error(),
		})
	}
	return c.JSON(game)
}