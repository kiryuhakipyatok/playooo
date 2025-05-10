package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"github.com/google/uuid"
	"time"
)

type CommentService interface {
	AddComment(ctx context.Context, whom, id, rid, body string) (*entities.Comment, error)
	GetComments(ctx context.Context, whose, id string, amount, page int) ([]entities.Comment, error)
}

type commentService struct{
	CommentRepository repositories.CommentRepository
	UserRepository repositories.UserRepository
	EventRepository repositories.EventRepository
	NewsRepository repositories.NewsRepository
	Transactor        repositories.Transactor
}

func NewCommentRepository(cr repositories.CommentRepository,ur repositories.UserRepository, er repositories.EventRepository, nr repositories.NewsRepository,t repositories.Transactor) CommentService{
	return &commentService{
		CommentRepository: cr,
		UserRepository: ur,
		EventRepository: er,
		NewsRepository: nr,
		Transactor: t,
	}
}

func(cs *commentService) AddComment(ctx context.Context,whom, id, rid, body string) (*entities.Comment, error){
	res,err:= cs.Transactor.WithinTransaction(ctx, func(c context.Context) (any, error) {
		user,err:=cs.UserRepository.FindById(c,id)
		if err!=nil{
			return nil,err
		}
		comment:=entities.Comment{
			Id: uuid.New(),
			AuthorId: user.Id,
			Body: body,
			Time: time.Now(),
		}
		if err:=cs.CommentRepository.Create(c,comment);err!=nil{
			return nil,err
		}
		if err:=cs.CommentRepository.Add(c,whom,rid,comment.Id.String());err!=nil{
			return nil,err
		}
		return &comment,nil
	})
	if err!=nil{
		return nil,err
	}
	return res.(*entities.Comment),nil
}

func(cs *commentService) GetComments(ctx context.Context, whose, id string, amount, page int) ([]entities.Comment, error){
	comments,err:=cs.CommentRepository.Fetch(ctx,whose,id,amount,page)
	if err!=nil{
		return nil,err
	}
	return comments,nil
}

