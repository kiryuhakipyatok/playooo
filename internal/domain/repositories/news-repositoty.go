package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"crap/internal/domain/entities"
)


type NewsRepository interface {
	Create(ctx context.Context, news entities.News) error
	Save(ctx context.Context, news entities.News) error
	FindById(ctx context.Context, id string) (*entities.News, error)
	Fetch(ctx context.Context, amount,page int) ([]entities.News, error)
}

type newsRepository struct {
	DB *pgx.Conn
}

func NewNewsRepository(db *pgx.Conn) NewsRepository {
	return &newsRepository{
		DB: db,
	}
}

func (nr *newsRepository) Create(ctx context.Context, news entities.News) error {
	if _,err:=nr.DB.Exec(ctx,"INSERT INTO news (id,title,body,time,link,picture) values($1,$2,$3,$4,$5,$6)",news.Id,news.Title,news.Body,news.Time,news.Link,news.Picture);err!=nil{
		return err
	}
	return nil
}

func (nr *newsRepository) Save(ctx context.Context, news entities.News) error {
	if _,err:=nr.DB.Exec(ctx,"UPDATE news SET title=$1,body=$2,time=$3,link=$4,picture=$5 WHERE id =$6)",news.Title,news.Body,news.Time,news.Link,news.Picture,news.Id);err!=nil{
		return err
	}
	return nil
}

func (nr *newsRepository) FindById(ctx context.Context, id string) (*entities.News, error) {
	news := entities.News{}
	if err:=nr.DB.QueryRow(ctx,"SELECT * FROM news WHERE id = $1",id).Scan(&news.Id,&news.Title,&news.Body,&news.Time,&news.Link,&news.Picture);err!=nil{
		return nil,err
	}
	return &news, nil
}

func (nr *newsRepository) Fetch(ctx context.Context, amount, page int) ([]entities.News, error) {
	somenews := []entities.News{}
	rows,err:=nr.DB.Query(ctx,"SELECT * FROM news OFFSET $1 LIMIT $2",page*amount-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		news:=entities.News{}
		if err:=rows.Scan(&news.Id,&news.Title,&news.Body,&news.Time,&news.Link,&news.Picture);err!=nil{
			return nil,err
		}
		somenews = append(somenews, news)
	}
	return somenews, nil
}