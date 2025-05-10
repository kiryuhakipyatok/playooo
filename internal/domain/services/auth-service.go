package services

import (
	"context"
	"crap/internal/config"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(c context.Context, login, tg, password string) (*entities.User, error)
	Login(c context.Context, login, password string) (*string, error)
	Profile(c context.Context, claims string) (*entities.User, error)
}

type authService struct {
	UserRepository repositories.UserRepository
	Config config.Config
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		UserRepository: userRepository,
	}
}

func (as *authService) Register(c context.Context, login, tg, password string) (*entities.User, error) {
	b,err:=as.UserRepository.ExistByLoginOrTg(c, login, tg) 
	if err!=nil{
		return nil,err
	}
	if b{
		return nil, errors.New("user with same login or telegram alredy exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	user := entities.User{
		Id:       uuid.New(),
		Login:    login,
		Telegram: tg,
		Password: hashPassword,
	}
	if err := as.UserRepository.Create(c, user); err != nil {
		return nil, err
	}
	return &user, nil
}

func(as *authService) Login(c context.Context, login,password string) (*string,error){
	if as.Config.Auth.Secret == "" {
		return nil,errors.New("error secret .env value is empty")
	}
	user, err := as.UserRepository.FindBy(c,"login", login)
	if err != nil {
		return nil,err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil,err
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.Id.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})
	token, err := claims.SignedString([]byte(as.Config.Auth.Secret))
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (as *authService) Profile(c context.Context, claims string) (*entities.User, error) {
	user, err := as.UserRepository.FindById(c, claims)
	if err != nil {
		return nil, err
	}
	return user, nil
}
