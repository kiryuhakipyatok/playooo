package handlers

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"errors"
	"strings"
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
func (gh *GamesHandler) AddGame(c *fiber.Ctx) error {
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
	gh.Logger.Infof("game %s added to user %v", request.Game,request.UserId)
	return c.JSON(fiber.Map{
		"message": "success",
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
func (gh *GamesHandler) DeleteGame(c *fiber.Ctx) error {
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
	if err := gh.GameService.DeleteGame(ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete game: " + err.Error(),
		})
	}
	gh.Logger.Infof("game %s deleted from user %v", request.Game,request.UserId)
	return c.JSON(fiber.Map{
		"message": "success",
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
func (gh *GamesHandler) GetGames(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-games")
	params := dto.PaginationRequest{}
	if err := c.QueryParser(&params); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := gh.Validator.Struct(params); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	games, err := gh.GameService.FetchGames(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get games: " + err.Error(),
		})
	}
	gh.Logger.Infof("games received: %v", params.Amount)
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
func (gh *GamesHandler) GetGame(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-game")
	title := c.Params("game")
	game, err := gh.GameService.GetByName(ctx, title)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get game: " + err.Error(),
		})
	}
	gh.Logger.Infof("game received: %v", title)
	return c.JSON(game)
}

// GetFilteredGames godoc
// @Summary Get filtered games
// @Description Get filtered list of games by params
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.GamesFilterRequest true "Filter parameters"
// @Success 200 {array} entities.Game "List of games"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games/filter [get]
func (gh *GamesHandler) GetFilteredGames(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-filtered-games")
	params := dto.GamesFilterRequest{}
	if err := c.QueryParser(&params); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := gh.Validator.Struct(params); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	gh.Logger.Infof("params: %+v",params)
	games, err := gh.GameService.GetFiltered(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to fetch filtered games: " + err.Error(),
		})
	}
	gh.Logger.Infof("filtered games received: %v", params.Amount)
	return c.JSON(games)
}

// GetSortedGames godoc
// @Summary Get sorted games
// @Description Get sorted list of games by params
// @Tags games
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.GamesSortRequest true "Sort parameters"
// @Success 200 {array} entities.Game "List of games"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /games/sort [get]
func (gh *GamesHandler) GetSortedGames(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, gh.Logger, "get-sorted-games")
	params := dto.GamesSortRequest{}
	if err := c.QueryParser(&params); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := gh.Validator.Struct(params); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	params.Direction = strings.ToUpper(params.Direction)
	if params.Direction == "" || params.Direction != "DESC" {
		params.Direction = "ASC"
	}
	games, err := gh.GameService.GetSorted(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to fetch sorted games: " + err.Error(),
		})
	}
	gh.Logger.Infof("sorted games received: %v", params.Amount)
	return c.JSON(games)
}
