package services

import (
	"context"
	"crap/internal/config"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type UserService interface {
	GetById(ctx context.Context, id string) (*entities.User, error)
	Fetch(ctx context.Context, amount, page int) ([]entities.User, error)
	UploadAvatar(ctx context.Context, id string, picture *multipart.FileHeader) error
	DeleteAvatar(ctx context.Context, id string) error
	RecordDiscord(ctx context.Context, id, ds string) error
	EditRating(ctx context.Context, id string, stars int) error
}

type userService struct {
	UserRepository repositories.UserRepository
	Transactor     repositories.Transactor
	Config config.Config
}

func NewUserService(ur repositories.UserRepository, t repositories.Transactor) UserService {
	return &userService{
		UserRepository: ur,
		Transactor:     t,
	}
}

func (us *userService) GetById(ctx context.Context, id string) (*entities.User, error) {
	user, err := us.UserRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) Fetch(ctx context.Context, amount, page int) ([]entities.User, error) {
	users, err := us.UserRepository.Fetch(ctx, amount, page)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) UploadAvatar(ctx context.Context, id string, picture *multipart.FileHeader) error {
	_, err := us.Transactor.WithinTransaction(ctx, func(c context.Context) (any, error) {
		user, err := us.UserRepository.FindById(c, id)
		if err != nil {
			return nil, err
		}
		uploadDir := "../../files/avatars"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, err
		}
		if _, err := os.Stat(uploadDir); err != nil {
			return nil, err
		}
		fileName := fmt.Sprintf("%s%s", user.Id, filepath.Ext(picture.Filename))
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
			host = us.Config.Server.Host
			port = us.Config.Server.Port
		)
		fileURL := fmt.Sprintf("http://%s:%s/files/avatars/%s", host, port, fileName)

		user.Avatar = fileURL
		if err := us.UserRepository.Save(c, *user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) DeleteAvatar(ctx context.Context, id string) error {
	user, err := us.GetById(ctx, id)
	if err != nil {
		return err
	}
	if user.Avatar == "" {
		return errors.New("user does not have an avatar")
	}
	var (
		host = us.Config.Server.Host
		port = us.Config.Server.Port
	)
	file := strings.TrimPrefix(user.Avatar, fmt.Sprintf("http://%s:%s/", host, port))
	if err := os.Remove(fmt.Sprintf("../../%s", file)); err != nil {
		return err
	}
	user.Avatar = ""
	if err := us.UserRepository.Save(ctx, *user); err != nil {
		return err
	}
	return nil
}

func (us *userService) RecordDiscord(ctx context.Context, id, ds string) error {
	user, err := us.GetById(ctx, id)
	if err != nil {
		return err
	}
	user.Discord = ds
	if err := us.UserRepository.Save(ctx, *user); err != nil {
		return err
	}
	return nil
}

func (us *userService) EditRating(ctx context.Context, id string, stars int) error {
	user, err := us.GetById(ctx, id)
	if err != nil {
		return err
	}
	user.NumberOfRatings++
	user.TotalRating += stars
	averageRating := float64(user.TotalRating) / float64(user.NumberOfRatings)
	user.Rating = math.Round(averageRating*2) / 2
	if err := us.UserRepository.Save(ctx, *user); err != nil {
		return err
	}
	return nil
}
