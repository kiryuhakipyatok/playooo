package handlers

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"time"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CommentsHandler struct {
	CommentService services.CommentService
	Logger         *logrus.Logger
	Validator      *validator.Validate
}

func NewCommentHandler(cs services.CommentService, l *logrus.Logger, v *validator.Validate) CommentsHandler {
	return CommentsHandler{
		CommentService: cs,
		Logger:         l,
		Validator:      v,
	}
}

// AddComment godoc
// @Summary Adding a comment
// @Description Creates a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AddCommentRequest true "Comment data"
// @Success 200 {object} entities.Comment 
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /comments [post]
func (ch *CommentsHandler) AddComment(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ch.Logger,"add-comment")
	request:=dto.AddCommentRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=ch.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	comment,err:=ch.CommentService.AddComment(ctx,request)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to add comment: " + err.Error(),
		})
	}
	ch.Logger.Infof("comment added: %s",comment.Id)
	return c.JSON(comment)
}

// GetComments godoc
// @Summary Receiving comments
// @Description Returns a list of comments by parameters
// @Tags comments
// @Produce json
// @Param whose query string false "whose"
// @Param userId query string false "user ID"
// @Param newsId query string false "news ID"
// @Param amount query int false "amount"
// @Param page query int false "page"
// @Success 200 {array} entities.Comment "List of comments"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /comments [get]
func (ch *CommentsHandler) GetComments(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ch.Logger,"get-comments")
	params:=dto.GetCommentsRequest{}
	if err:=c.QueryParser(&params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	if err:=ch.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	comments,err:=ch.CommentService.GetComments(ctx,params)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get comments: " + err.Error(),
		})
	}
	ch.Logger.Infof("comments recieved: %s",params.Id)
	return c.JSON(comments)
}

