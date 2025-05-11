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

func (eh *EventsHandler) CreateEvent(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"create-event")
	request:=dto.CreateEventRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=eh.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	event,err:=eh.EventService.CreateEvent(ctx,request)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to create event: " +  err.Error(),
		})
	}
	responce:=dto.EventResponse{
		Id: event.Id,
		AuthorId: event.AuthorId,
		Time: event.Time,
	}
	return c.JSON(responce)
}

func(eh *EventsHandler) GetEvent(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"get-event-by-id")
	id:=c.Params("id")
	event,err:=eh.EventService.GetById(ctx,id)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get event: "+err.Error(),
		})
	}
	return c.JSON(event)
}

func(eh *EventsHandler) GetEvents(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"get-events")
	params:=dto.PaginationRequest{}
	if err:=c.QueryParser(params);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=eh.Validator.Struct(params);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	events,err:=eh.EventService.FetchEvents(ctx,params)
	if err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to get events: "+err.Error(),
		})
	}
	return c.JSON(events)
}

func(eh *EventsHandler) DeleteEvent(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"delete-event")
	id:=c.Params("id")
	if err:=eh.EventService.DeleteEvent(ctx,id);err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to delete event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(eh *EventsHandler) Join(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"join-to-event")
	request:=dto.JoinToEventRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=eh.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	if err:=eh.EventService.Join(ctx,request);err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to join to event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(eh *EventsHandler) Unjoin(c *fiber.Ctx) error{
	ctx,cancel:=context.WithTimeout(c.Context(),time.Second*5)
	defer cancel()
	eH:=errh.NewErrorHander(c,eh.Logger,"unjoin-from-event")
	request:=dto.UnjoinFromEventRequest{}
	if err:=c.BodyParser(&request);err!=nil{
		return errh.ParseRequestError(eH,err)
	}
	if err:=eh.Validator.Struct(request);err!=nil{
		return errh.ValidateRequestError(eH,err)
	}
	if err:=eh.EventService.Unjoin(ctx,request);err!=nil{
		if errors.Is(err,context.DeadlineExceeded){
			return errh.RequestTimedOut(eH,err)
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to unjoin from event: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

