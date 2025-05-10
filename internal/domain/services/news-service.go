package services

import (
	"context"
	"crap/internal/config"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type NewsService interface {
	CreateNews(ctx context.Context, title, body, link string, picture *multipart.FileHeader) (*entities.News, error)
	GetById(ctx context.Context, id string) (*entities.News, error)
	FetchNews(ctx context.Context, amount, page int) ([]entities.News, error)
}

type newsService struct {
	NewsRepository    repositories.NewsRepository
	Transactor repositories.Transactor
	Config config.Config
}

func NewNewsService(nr repositories.NewsRepository, t repositories.Transactor) NewsService {
	return &newsService{
		NewsRepository: nr,
		Transactor: t,
	}
}

func (ns *newsService) CreateNews(ctx context.Context, title, body, link string, picture *multipart.FileHeader) (*entities.News, error){
	res,err:=ns.Transactor.WithinTransaction(ctx,func(c context.Context) (any, error) {
			news := entities.News{
			Id:    uuid.New(),
			Title: title,
			Body:  body,
			Time:  time.Now(),
			Link:  link,
		}
	
		uploadDir := "../../files/news-pictures"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, err
		}
		if _, err := os.Stat(uploadDir);err!=nil {
			return nil, err
		}
		if filepath.Ext(picture.Filename) != ".png" || filepath.Ext(picture.Filename) != ".jpg"{
			return nil,errors.New("incorrect picture format")
		}
		fileName := fmt.Sprintf("%s-news-picture%s", news.Id, filepath.Ext(picture.Filename))
		filepath := filepath.Join(uploadDir, fileName)
		dst, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()
		src, err := picture.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			return nil, err
		}

		var (
			host = ns.Config.Server.Host
			port = ns.Config.Server.Port
		)

		fileURL := fmt.Sprintf("http://%s:%s/files/news-pictures/%s", host, port, fileName)

		news.Picture = fileURL
		if err := ns.NewsRepository.Create(c, news); err != nil {
			return nil, err
		}
		return nil,nil
	})
	if err!=nil{
		return nil,err
	}

	return res.(*entities.News), nil
}

func (ns *newsService) GetById(ctx context.Context, id string) (*entities.News, error){
	news,err:=ns.NewsRepository.FindById(ctx,id)
	if err!=nil{
		return nil,err
	}
	return news,nil
}

func (ns *newsService) FetchNews(ctx context.Context, amount, page int) ([]entities.News, error){
	news,err:=ns.NewsRepository.Fetch(ctx,amount,page)
	if err!=nil{
		return nil,err
	}
	return news,nil
}

