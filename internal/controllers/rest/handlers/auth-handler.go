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
	Config 		*config.Config
}

func NewAuthHandler(as services.AuthService, l *logrus.Logger, v *validator.Validate, cfg *config.Config) AuthHandler{
	return AuthHandler{
		AuthService: as,
		Validator: v,
		Logger: l,
		Config: cfg,
	}
}

// Register godoc
// @Summary User registration
// @Description Creates a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 200 {object}  dto.RegisterRequest
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /auth/register [post]
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
	user,err:=ah.AuthService.Register(ctx,request)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "register failed: " + err.Error(),
		})
	}
	ah.Logger.Infof("user registered: %v",user.Id)
	response:=dto.RegisterResponse{
		Id: user.Id,
		Login: user.Login,
		Telegram: user.Telegram,
		Date: user.DateOfRegister,
	}
	return c.JSON(response)
}

// Login godoc
// @Summary User authentication
// @Description User login and receiving JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login data"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /auth/login [post]
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
	token,err:=ah.AuthService.Login(ctx,request)
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

// Logout godoc
// @Summary Logout
// @Description Clears the JWT cookie
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /auth/logout [post]
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

// Profile godoc
// @Summary Getting a logged in profile
// @Description Returns the profile data of the current user
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entities.User
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /auth/profile [get]
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