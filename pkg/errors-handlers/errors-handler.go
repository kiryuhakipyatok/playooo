package error_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ErrorHandler struct {
	Ctx         *fiber.Ctx
	Logger      *logrus.Logger
	RequestType string
}

func NewErrorHander(c *fiber.Ctx, l *logrus.Logger, rt string) ErrorHandler {
	return ErrorHandler{
		Ctx:         c,
		Logger:      l,
		RequestType: rt,
	}
}

func ParseRequestError(eh ErrorHandler, err error) error {
	eh.Ctx.Status(fiber.StatusBadRequest)
	eh.Logger.WithError(err).Errorf("failed to parse %s request", eh.RequestType)
	return eh.Ctx.JSON(fiber.Map{
		"error": "bad request: " + err.Error(),
	})
}

func ValidateRequestError(eh ErrorHandler, err error) error {
	eh.Ctx.Status(fiber.StatusBadRequest)
	eh.Logger.WithError(err).Errorf("%s request validation failed", eh.RequestType)
	return eh.Ctx.JSON(fiber.Map{
		"error": "validation failed: " + err.Error(),
	})
}

func RequestTimedOut(eh ErrorHandler, err error) error {
	eh.Ctx.Status(fiber.StatusRequestTimeout)
	eh.Logger.Debugf("%s request timed out", eh.RequestType)
	return eh.Ctx.JSON(fiber.Map{
		"error": "request timed out: " + err.Error(),
	})
}
