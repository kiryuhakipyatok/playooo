package repositories

import (
	"context"
	"crap/internal/domain/entities"
	"github.com/jackc/pgx/v5"
)

type CommentRepository interface {
	Create(ctx context.Context, comment entities.Comment) error
	FetchFromUser(ctx context.Context,id string, amount, page int) ([]entities.Comment, error)
	FetchFromEvent(ctx context.Context, id string, amount, page int) ([]entities.Comment, error)
	FetchFromNews(ctx context.Context,id string, amount, page int) ([]entities.Comment, error)
	AddToUser(ctx context.Context,id, cid string) error
	AddToNews(ctx context.Context,id, cid string) error
	AddToEvent(ctx context.Context,id, cid string) error
}


type commentRepository struct {
	DB *pgx.Conn
}

func NewCommentRepository(db *pgx.Conn) CommentRepository {
	return &commentRepository{
		DB: db,
	}
}

func (cr *commentRepository) Create(ctx context.Context, comment entities.Comment) error {
	if _,err:=cr.DB.Exec(ctx,"INSERT INTO comments (id,author_id,body,time) values ($1,$2,$3,$4)", comment.Id,comment.AuthorId,comment.Body,comment.Time);err!=nil{
		return err
	}
	return nil	
}

func (cr *commentRepository) Fetch(ctx context.Context,whose, id string, amount, page int) ([]entities.Comment, error) {
	comments := []entities.Comment{}
	var query string
	switch whose{
	case "user":
		query = "SELECT * FROM comments c JOIN users_comments uc ON c.id=uc.comment_id WHERE uc.user_id = $1 ORDER BY time OFFSET $2 LIMIT $3"
	case "event":
		query = "SELECT * FROM comments c JOIN events_comments ec ON c.id=ec.comment_id WHERE ec.event_id = $1 ORDER BY time OFFSET $2 LIMIT $3"
	case "news":
		query = "SELECT * FROM comments c JOIN news_comments uc ON c.id=nc.comment_id WHERE nc.news_id = $1 ORDER BY time OFFSET $2 LIMIT $3"
	}
	rows,err:=cr.DB.Query(ctx,query,id,page*amount-amount, amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		comment:=entities.Comment{}
		if err:=rows.Scan(&comment.Id,&comment.AuthorId,&comment.Body,&comment.Time);err!=nil{
			return nil,err
		}
		comments=append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) FetchFromUser(ctx context.Context,id string, amount, page int) ([]entities.Comment, error){
	comments := []entities.Comment{}
	query := "SELECT * FROM comments c JOIN users_comments uc ON c.id=uc.comment_id WHERE uc.user_id = $1 OFFSET $2 LIMIT $3"
	rows,err:=cr.DB.Query(ctx,query,id,page*amount-amount, amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		comment:=entities.Comment{}
		if err:=rows.Scan(&comment.Id,&comment.AuthorId,&comment.Body,&comment.Time);err!=nil{
			return nil,err
		}
		comments=append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) FetchFromEvent(ctx context.Context, id string, amount, page int) ([]entities.Comment, error){
	comments := []entities.Comment{}
	query := "SELECT * FROM comments c JOIN events_comments ec ON c.id=ec.comment_id WHERE ec.event_id = $1 OFFSET $2 LIMIT $3"
	rows,err:=cr.DB.Query(ctx,query,id,page*amount-amount, amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		comment:=entities.Comment{}
		if err:=rows.Scan(&comment.Id,&comment.AuthorId,&comment.Body,&comment.Time);err!=nil{
			return nil,err
		}
		comments=append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) FetchFromNews(ctx context.Context,id string, amount, page int) ([]entities.Comment, error){
	comments := []entities.Comment{}
	query := "SELECT * FROM comments c JOIN news_comments nc ON c.id=nc.comment_id WHERE nc.news_id = $1 OFFSET $2 LIMIT $3"
	rows,err:=cr.DB.Query(ctx,query,id,page*amount-amount, amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		comment:=entities.Comment{}
		if err:=rows.Scan(&comment.Id,&comment.AuthorId,&comment.Body,&comment.Time);err!=nil{
			return nil,err
		}
		comments=append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) AddToUser(ctx context.Context,id, cid string) error{
	query := "INSERT INTO users_comments (comment_id,user_id) values($1,$2)"
	if _,err:=cr.DB.Exec(ctx,query,cid,id);err!=nil{
		return err
	}
	return nil
}

func (cr *commentRepository) AddToEvent(ctx context.Context,id, cid string) error{
	query := "INSERT INTO events_comments (comment_id,event_id) values($1,$2)"
	if _,err:=cr.DB.Exec(ctx,query,cid,id);err!=nil{
		return err
	}
	return nil
}

func (cr *commentRepository) AddToNews(ctx context.Context,id, cid string) error{
	query := "INSERT INTO news_comments (comment_id,news_id) values($1,$2)"
	if _,err:=cr.DB.Exec(ctx,query,cid,id);err!=nil{
		return err
	}
	return nil
}

