package handlers

import (
	"context"
	"crap/internal/config"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
) 

type AuthHandler struct {
	AuthService services.AuthService
	Validator    *validator.Validate
	Logger       *logrus.Logger
	Config 		config.Config
}

func NewAuthHandler(as services.AuthService, l *logrus.Logger, v *validator.Validate) AuthHandler{
	return AuthHandler{
		AuthService: as,
		Validator: v,
		Logger: l,
	}
}

func(ah *AuthHandler) Register(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ah.Logger,"register")
	request:=dto.RegisterRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=ah.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	user,err:=ah.AuthService.Register(ctx,request.Login,request.Telegram,request.Password)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "register failed: " + err.Error(),
		})
	}
	response:=dto.RegisterResponse{
		Id: user.Id,
		Login: user.Login,
		Telegram: user.Telegram,
	}
	return c.JSON(response)
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ah.Logger,"login")
	request:=dto.LoginRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=ah.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	token,err:=ah.AuthService.Login(ctx,request.Login,request.Password)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "login failed: " + err.Error(),
		})
	}
	jwt:=fiber.Cookie{
		Name: "jwt",
		Value: *token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(&jwt)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(ah *AuthHandler) Logout(c *fiber.Ctx) error{
	jwt:=fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(&jwt)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func(ah *AuthHandler) Profile(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,ah.Logger,"profile")
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(ah.Config.Auth.Secret), nil
	})
	if err!=nil{
		ah.Logger.WithError(err).Info("unauthenticated")
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	user, err := ah.AuthService.Profile(ctx, claims.Issuer)
	if err != nil {
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		ah.Logger.WithError(err).Info("unauthenticated")
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "not found",
		})
	}
	return c.JSON(user)
}