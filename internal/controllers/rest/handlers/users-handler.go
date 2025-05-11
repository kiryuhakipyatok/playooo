package handlers

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserService services.UserService
	Logger      *logrus.Logger
	Validator   *validator.Validate
}

func NewUserHandler(us services.UserService, l *logrus.Logger, v *validator.Validate) UserHandler {
	return UserHandler{
		UserService: us,
		Logger:      l,
		Validator:   v,
	}
}

func(uh *UserHandler) GetById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "get-user")
	id := c.Params("id")
	user, err := uh.UserService.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get user: " + err.Error(),
		})
	}
	return c.JSON(user)
}

func(uh *UserHandler) GetUsers(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "get-user")
	params:=dto.PaginationRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=uh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	users,err:=uh.UserService.Fetch(ctx,params.Amount,params.Page)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get users: " + err.Error(),
		})
	}
	return c.JSON(users)
}

func(uh *UserHandler) UploadAvatar(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "upload-avatar")
	request := dto.UploadAvatarRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := uh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=uh.UserService.UploadAvatar(ctx,request.UserId,request.Picture);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to upload avatar: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(uh *UserHandler) RecordDiscord(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "record-discord")
	request := dto.RecordDiscordRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := uh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=uh.UserService.RecordDiscord(ctx,request.UserId,request.Discord);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to record discord: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"succes",
	})
}

func(uh *UserHandler) DeleteAvatar(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "upload-avatar")
	id:=c.Params("id")
	if err:=uh.UserService.DeleteAvatar(ctx,id);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete avatar: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(uh *UserHandler) EditRating(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "edit-rating")
	request := dto.EditRatingRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := uh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=uh.UserService.EditRating(ctx,request.UserId,request.Stars);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to edit rating: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"succes",
	})
}