package handlers

import (
	"context"
	"crap/internal/domain/services"
	"crap/internal/dto"
	errh "crap/pkg/errors-handlers"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EventsHandler struct {
	EventService services.EventService
	Logger       *logrus.Logger
	Validator    *validator.Validate
}

func NewEventHandler(es services.EventService, l *logrus.Logger, v *validator.Validate) EventsHandler {
	return EventsHandler{
		EventService: es,
		Logger:       l,
		Validator:    v,
	}
}

// CreateEvent godoc
// @Summary Create an event
// @Description Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateEventRequest true "Event data"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events [post]
func (eh *EventsHandler) CreateEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "create-event")
	request := dto.CreateEventRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	event, err := eh.EventService.CreateEvent(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to create event: " + err.Error(),
		})
	}
	// responce:=dto.EventResponse{
	// 	Id: event.Id,
	// 	AuthorId: event.AuthorId,
	// 	Time: event.Time,
	// }
	return c.JSON(event)
}

// GetEvent godoc
// @Summary Getting event by ID
// @Description Returns an event by its id
// @Tags events
// @Produce json
// @Param id path string true "event Id"
// @Success 200 {object} entities.Event
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/{id} [get]
func (eh *EventsHandler) GetEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "get-event-by-id")
	id := c.Params("id")
	event, err := eh.EventService.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get event: " + err.Error(),
		})
	}
	return c.JSON(event)
}

// GetEvents godoc
// @Summary Getting a list of events
// @Description Returns a list of events with pagination
// @Tags events
// @Produce json
// @Param amount query int false "amount"
// @Param page query int false "page"
// @Success 200 {array} entities.Event
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events [get]
func (eh *EventsHandler) GetEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "get-events")
	params := dto.PaginationRequest{}
	if err := c.QueryParser(&params); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(params); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	events, err := eh.EventService.FetchEvents(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to get events: " + err.Error(),
		})
	}
	return c.JSON(events)
}

// DeleteEvent godoc
// @Summary Deleting an event
// @Description Deletes an event by its ID
// @Tags events
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Event ID"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/{id} [delete]
func (eh *EventsHandler) DeleteEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "delete-event")
	id := c.Params("id")
	if err := eh.EventService.DeleteEvent(ctx, id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to delete event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Join godoc
// @Summary Joining the event
// @Description Adds a user to the event participants
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.JoinToEventRequest true "Data for join to event"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/join [post]
func (eh *EventsHandler) Join(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "join-to-event")
	request := dto.JoinToEventRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err := eh.EventService.Join(ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to join to event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Unjoin godoc
// @Summary Exit event
// @Description Removes a user from the event participants
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.UnjoinFromEventRequest true "Data for unjoin from event"
// @Success 200 {object} object "{\"message\":\"string\"}"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/unjoin [post]
func (eh *EventsHandler) Unjoin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "unjoin-from-event")
	request := dto.UnjoinFromEventRequest{}
	if err := c.BodyParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	if err := eh.EventService.Unjoin(ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to unjoin from event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// GetFilteredEvents godoc
// @Summary Get filtered events
// @Description Get filtered list of events by params
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.EventsFilterRequest true "Filter parameters"
// @Success 200 {array} entities.Event "List of events"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/filter [get]
func (eh *EventsHandler) GetFilteredEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "get-filtered-events")
	request := dto.EventsFilterRequest{}
	if err := c.QueryParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	events, err := eh.EventService.GetFiltered(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to fetch filtered events: " + err.Error(),
		})
	}
	return c.JSON(events)
}

// GetSortedEvents godoc
// @Summary Get sorted events
// @Description Get sorted list of events by params
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request query dto.EventsSortRequest true "Sort parameters"
// @Success 200 {array} entities.Event "List of events"
// @Failure 400 {object} object "{\"error\":\"string\"}"
// @Failure 408 {object} object "{\"error\":\"string\"}"
// @Failure 500 {object} object "{\"error\":\"string\"}"
// @Router /events/sort [get]
func (eh *EventsHandler) GetSortedEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()
	eH := errh.NewErrorHander(c, eh.Logger, "get-sorted-events")
	request := dto.EventsSortRequest{}
	if err := c.QueryParser(&request); err != nil {
		return errh.ParseRequestError(eH, err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return errh.ValidateRequestError(eH, err)
	}
	request.Direction = strings.ToUpper(request.Direction)
	if request.Direction == "" || request.Direction != "DESC" {
		request.Direction = "ASC"
	}
	events, err := eh.EventService.GetSorted(ctx, request)
	if err!=nil{
		if errors.Is(err, context.DeadlineExceeded) {
			return errh.RequestTimedOut(eH, err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to fetch sorted events: " + err.Error(),
		})
	}
	return c.JSON(events)
}
