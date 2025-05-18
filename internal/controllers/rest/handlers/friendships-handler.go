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

// AddFriend godoc
// @Summary Add a new friend
// @Description Send a friend request to another user
// @Tags friends
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AddFriendRequest true "Add Friend Request"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /friends/add [post]
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

// GetFriends godoc
// @Summary Get user's friends list
// @Description Retrieve the list of friends for a user
// @Tags friends
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.GetFriendsRequest true "Get Friends Request"
// @Success 200 {array} entities.User "List of friends"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /friends [get]
func (fh *FriendshipsHandler) GetFriends(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "get-friends")
	params:=dto.GetFriendsRequest{}
	if err:=c.QueryParser(&params);err!=nil{
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

// CancelFriendship godoc
// @Summary Cancel a friendship
// @Description Remove a friend or cancel a pending friend request
// @Tags friends
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CancelFriendshipRequest true "Cancel Friendship Request"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /friends/cancel [post]
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

// AcceptFriendship godoc
// @Summary Accept a friend request
// @Description Accept a pending friend request from another user
// @Tags friends
// @Accept json
// @Produce json
// @Param request body dto.AcceptFriendshipRequest true "Accept Friendship Request"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /friends/accept [post]
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

// GetFriendsRequests godoc
// @Summary Get friend requests
// @Description Retrieve pending friend requests for a user
// @Tags friends
// @Accept json
// @Produce json
// @Param request query dto.GetFriendsReqRequests true "Get Friend Requests"
// @Success 200 {array} entities.User "List of friend requests"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /friends/requests [get]
func (fh *FriendshipsHandler) GetFriendsRequests(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, fh.Logger, "get-friends-requests")
	params:=dto.GetFriendsReqRequests{}
	if err:=c.QueryParser(&params);err!=nil{
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