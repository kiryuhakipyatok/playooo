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

// DeleteNotification godoc
// @Summary Delete notification
// @Description Delete specific notification by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /notifications/{id} [delete]
func (nh *NotificationsHandler) DeleteNotification(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "delete-notification")
	id := c.Query("id")
	nid := c.Query("nid")
	if err := nh.NotificationService.DeleteNotification(ctx, id, nid); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete notification: " + err.Error(),
		})
	}
	nh.Logger.Infof("notification deleted: %v",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// DeleteAllNotifications godoc
// @Summary Delete all notifications
// @Description Delete all notifications for user
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /notifications/all/{id} [delete]
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
	nh.Logger.Infof("all notifications deleted: %v",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

// GetNotifications godoc
// @Summary Get notifications
// @Description Get paginated list of notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.GetNotificationsRequest true "Pagination and filter parameters"
// @Success 200 {array} entities.Notification "List of notifications"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /notifications [get]
func (nh *NotificationsHandler) GetNotifications(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, nh.Logger, "get-notifications")
	params := dto.GetNotificationsRequest{}
	if err:=c.QueryParser(&params);err!=nil{
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
	nh.Logger.Infof("notification received: %v",params.Amount)
	return c.JSON(notifications)
}


