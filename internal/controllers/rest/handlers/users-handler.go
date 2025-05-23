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

type UsersHandler struct {
	UserService services.UserService
	Logger      *logrus.Logger
	Validator   *validator.Validate
}

func NewUsersHandler(us services.UserService, l *logrus.Logger, v *validator.Validate) UsersHandler {
	return UsersHandler{
		UserService: us,
		Logger:      l,
		Validator:   v,
	}
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get detailed information about specific user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users/{id} [get]
func(uh *UsersHandler) GetUser(c *fiber.Ctx) error {
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
	uh.Logger.Infof("user received: %v", user.Id)
	return c.JSON(user)
}

// GetUsers godoc
// @Summary Get users list
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.PaginationRequest true "Pagination parameters"
// @Success 200 {array} entities.User "List of users"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users [get]
func(uh *UsersHandler) GetUsers(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "get-user")
	params:=dto.PaginationRequest{}
	if err:=c.QueryParser(&params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=uh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	users,err:=uh.UserService.Fetch(ctx,params)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get users: " + err.Error(),
		})
	}
	uh.Logger.Infof("users received: %v", params.Amount)
	return c.JSON(users)
}

// UploadAvatar godoc
// @Summary Upload user avatar
// @Description Upload or update user avatar image
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.UploadAvatarRequest true "Avatar upload data"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users/avatar [post]
func(uh *UsersHandler) UploadAvatar(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, uh.Logger, "upload-avatar")
	request := dto.UploadAvatarRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	id:=c.FormValue("id")
	picture,err:=c.FormFile("picture")
	if err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "failed to get picture: " + err.Error(),
		})
	}
	request.UserId=id
	request.Picture=picture
	if err := uh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err:=uh.UserService.UploadAvatar(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to upload avatar: " + err.Error(),
		})
	}
	uh.Logger.Infof("avatar uploaded: %v",request.UserId)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// RecordDiscord godoc
// @Summary Record Discord association
// @Description Link user account with Discord
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.RecordDiscordRequest true "Discord association data"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users/discord [post]
func(uh *UsersHandler) RecordDiscord(c *fiber.Ctx) error{
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
	if err:=uh.UserService.RecordDiscord(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to record discord: " + err.Error(),
		})
	}
	uh.Logger.Infof("discord recorded: %v", request.UserId)
	return c.JSON(fiber.Map{
		"message":"succes",
	})
}

// DeleteAvatar godoc
// @Summary Delete user avatar
// @Description Remove user's avatar image
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users/avatar/{id} [delete]
func(uh *UsersHandler) DeleteAvatar(c *fiber.Ctx) error{
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
	uh.Logger.Infof("avatar deleted: %v", id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// EditRating godoc
// @Summary Edit user rating
// @Description Update user's rating value
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.EditRatingRequest true "Rating edit data"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /users/rating [patch]
func(uh *UsersHandler) EditRating(c *fiber.Ctx) error{
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
	if err:=uh.UserService.EditRating(ctx,request);err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to edit rating: " + err.Error(),
		})
	}
	uh.Logger.Infof("rating edited: %v", request.UserId)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}