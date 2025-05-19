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

type NewsHandler struct {
	NewsService services.NewsService
	Logger      *logrus.Logger
	Validator   *validator.Validate
}

func NewNewsHandler(fs services.NewsService, l *logrus.Logger, v *validator.Validate) NewsHandler {
	return NewsHandler{
		NewsService: fs,
		Logger:      l,
		Validator:   v,
	}
}

// CreateNews godoc
// @Summary Create news article
// @Description Create a new news article
// @Tags news
// @Accept json
// @Produce json
// @Param request body dto.CreateNewsRequest true "News creation data"
// @Success 200 {object} dto.NewsResponse "News data"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /news [post]
func(nh *NewsHandler) CreateNews(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "create-news")
	request := dto.CreateNewsRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := nh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	news, err := nh.NewsService.CreateNews(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to create news: " + err.Error(),
		})
	}
	responce:=dto.NewsResponse{
		Id: news.Id,
		Title: news.Title,
	}
	nh.Logger.Infof("news created: %v", news.Id)
	return c.JSON(&responce)
}

// GetNews godoc
// @Summary Get news by ID (path)
// @Description Get news article by its ID (passed as path parameter)
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "News ID"
// @Success 200 {object} entities.News "News data"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /news/{id} [get]
func(nh *NewsHandler) GetNews(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "get-news")
	id:=c.Params("id")
	news,err:=nh.NewsService.GetById(ctx,id)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get news: "+ err.Error(),
		})
	}
	nh.Logger.Infof("news received: %v", news.Id)
	return c.JSON(news)
}

// GetSomeNews godoc
// @Summary Get paginated news
// @Description Get paginated list of news articles
// @Tags news
// @Accept json
// @Produce json
// @Param request query dto.PaginationRequest true "Pagination parameters"
// @Success 200 {array} entities.News "List of news"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /news [get]
func(nh *NewsHandler) GetSomeNews(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "get-news")
	params:=dto.PaginationRequest{}
	if err:=c.QueryParser(&params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=nh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	someNews,err:=nh.NewsService.FetchNews(ctx,params)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get some news: "+ err.Error(),
		})
	}
	nh.Logger.Infof("some news received: %v", params.Amount)
	return c.JSON(someNews)
}