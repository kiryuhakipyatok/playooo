package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"crap/internal/dto"
	"time"

	"github.com/google/uuid"
)

type CommentService interface {
	AddComment(ctx context.Context, req dto.AddCommentRequest) (*entities.Comment, error)
	GetComments(ctx context.Context, req dto.GetCommentsRequest) ([]entities.Comment, error)
}

type commentService struct{
	CommentRepository repositories.CommentRepository
	UserRepository repositories.UserRepository
	EventRepository repositories.EventRepository
	NewsRepository repositories.NewsRepository
	Transactor        repositories.Transactor
}

func NewCommentService(cr repositories.CommentRepository,ur repositories.UserRepository, er repositories.EventRepository, nr repositories.NewsRepository,t repositories.Transactor) CommentService{
	return &commentService{
		CommentRepository: cr,
		UserRepository: ur,
		EventRepository: er,
		NewsRepository: nr,
		Transactor: t,
	}
}

func(cs *commentService) AddComment(ctx context.Context,req dto.AddCommentRequest) (*entities.Comment, error){
	res,err:= cs.Transactor.WithinTransaction(ctx, func(c context.Context) (any, error) {
		user,err:=cs.UserRepository.FindById(c,req.UserId)
		if err!=nil{
			return nil,err
		}
		comment:=entities.Comment{
			Id: uuid.New(),
			AuthorId: user.Id,
			Body: req.Body,
			Time: time.Now(),
		}
		if err:=cs.CommentRepository.Create(c,comment);err!=nil{
			return nil,err
		}
		switch req.Whom{
		case "users":
			receiver,err:=cs.UserRepository.FindById(c,req.ReceiverId)
			if err!=nil{
				return nil,err
			}
			if err:=cs.CommentRepository.AddToUser(c,receiver.Id.String(),comment.Id.String());err!=nil{
				return nil,err
			}
		case "events":
			event,err:=cs.EventRepository.FindById(c,req.ReceiverId)
			if err!=nil{
				return nil,err
			}
			if err:=cs.CommentRepository.AddToEvent(c,event.Id.String(),comment.Id.String());err!=nil{
				return nil,err
			}
		case "news":
			news,err:=cs.NewsRepository.FindById(c,req.ReceiverId)
			if err!=nil{
				return nil,err
			}
			if err:=cs.CommentRepository.AddToNews(c,news.Id.String(),comment.Id.String());err!=nil{
				return nil,err
			}
		}
		return &comment,nil
	})
	if err!=nil{
		return nil,err
	}
	return res.(*entities.Comment),nil
}

func(cs *commentService) GetComments(ctx context.Context, req dto.GetCommentsRequest) ([]entities.Comment, error){
	comments:=[]entities.Comment{}
	switch req.Whose{
	case "user":
		var err error
		comments,err=cs.CommentRepository.FetchFromUser(ctx,req.Id,req.Amount,req.Page)
		if err!=nil{
			return nil,err
		}
	case "event":
		var err error
		comments,err=cs.CommentRepository.FetchFromEvent(ctx,req.Id,req.Amount,req.Page)
		if err!=nil{
			return nil,err
		}
	case "news":
		var err error
		comments,err=cs.CommentRepository.FetchFromNews(ctx,req.Id,req.Amount,req.Page)
		if err!=nil{
			return nil,err
		}
	}
	return comments,nil
}

