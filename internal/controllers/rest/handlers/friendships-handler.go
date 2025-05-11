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

type FriendshipsHandler struct {
	FriendshipsService services.FriendshipsService
	Logger             *logrus.Logger
	Validator          *validator.Validate
}

func NewFriendshipsHandler(fs services.FriendshipsService, l *logrus.Logger, v *validator.Validate) FriendshipsHandler {
	return FriendshipsHandler{
		FriendshipsService: fs,
		Logger:             l,
		Validator:          v,
	}
}

func (fh *FriendshipsHandler) AddFriend(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "add-friend")
	request := dto.AddFriendRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := fh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err := fh.FriendshipsService.AddFriend(ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to add friend: "+ err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (fh *FriendshipsHandler) GetFriends(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "get-friends")
	params:=dto.GetFriendsRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=fh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	friends,err:=fh.FriendshipsService.GetFriends(ctx,params)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get friends: "+ err.Error(),
		})
	}
	return c.JSON(friends)
}

func (fh *FriendshipsHandler) CancelFriendship(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "cancel-friendship")
	request := dto.CancelFriendshipRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := fh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=fh.FriendshipsService.CancelFriendship(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to cancel friendship: "+ err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (fh *FriendshipsHandler) AcceptFriendship(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "accept-friendship")
	request := dto.AcceptFriendshipRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := fh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=fh.FriendshipsService.AcceptFriendship(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to accept friendship: "+ err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (fh *FriendshipsHandler) GetFriendsRequests(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "get-friends-requests")
	params:=dto.GetFriendsReqRequests{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=fh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	requests,err:=fh.FriendshipsService.GetFriendRequests(ctx,params)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get friend request: "+ err.Error(),
		})
	}
	return c.JSON(requests)
}