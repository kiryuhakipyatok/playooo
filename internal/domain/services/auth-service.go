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
	Register(ctx context.Context, login, tg, password string) (*entities.User, error)
	Login(ctx context.Context, login, password string) (*string, error)
	Profile(ctx context.Context, claims string) (*entities.User, error)
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

func (as *authService) Register(ctx context.Context, login, tg, password string) (*entities.User, error) {
	b,err:=as.UserRepository.ExistByLoginOrTg(ctx, login, tg) 
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
	if err := as.UserRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	return &user, nil
}

func(as *authService) Login(ctx context.Context, login,password string) (*string,error){
	if as.Config.Auth.Secret == "" {
		return nil,errors.New("error secret .env value is empty")
	}
	user, err := as.UserRepository.FindBy(ctx,"login", login)
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

func (as *authService) Profile(ctx context.Context, claims string) (*entities.User, error) {
	user, err := as.UserRepository.FindById(ctx, claims)
	if err != nil {
		return nil, err
	}
	return user, nil
}
