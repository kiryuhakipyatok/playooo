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
	comment,err:=ch.CommentService.AddComment(ctx,request.Whom,request.UserId,request.ReceiverId,request.Body)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to add comment: " + err.Error(),
		})
	}
	return c.JSON(comment)
}

func (ch *CommentsHandler) GetComments(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ch.Logger,"get-comments")
	params:=dto.GetCommentsRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	if err:=ch.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	comments,err:=ch.CommentService.GetComments(ctx,params.Whose,params.UserId,params.Amount,params.Page)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get comments: " + err.Error(),
		})
	}
	return c.JSON(comments)
}

