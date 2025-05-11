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

type NotificationsHandler struct {
	NotificationService services.NotificationService
	Logger              *logrus.Logger
	Validator           *validator.Validate
}

func NewNotificationsHandler(ns services.NotificationService, l *logrus.Logger, v *validator.Validate) NotificationsHandler {
	return NotificationsHandler{
		NotificationService: ns,
		Logger:              l,
		Validator:           v,
	}
}

func (nh *NotificationsHandler) DeleteNotification(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "delete-notification")
	id := c.Params("id")
	if err := nh.NotificationService.DeleteNotification(ctx, id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete notification: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (nh *NotificationsHandler) DeleteAllNotifications(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "delete-all-notifications")
	id := c.Params("id")
	if err := nh.NotificationService.DeleteAllNotifications(ctx, id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete all notifications: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (nh *NotificationsHandler) GetNotifications(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "get-notifications")
	params := dto.GetNotificationsRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err := nh.Validator.Struct(params); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	notifications,err := nh.NotificationService.FetchNotifications(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get notifications: " + err.Error(),
		})
	}
	return c.JSON(notifications)
}


