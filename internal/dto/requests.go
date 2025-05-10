package dto

import "mime/multipart"

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,max=100"`
	Telegram string `json:"telegram" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,max=100"`
	Password string `json:"password" validate:"required"`
}

type UserCommentRequest struct {
	AuthorId   string `json:"author-id" validate:"required"`
	ReceiverId string `json:"reciever-id" validate:"required"`
	Body       string `json:"body" validate:"max=150"`
}

type NewsCommentRequest struct {
	AuthorId   string `json:"author-id" validate:"required"`
	ReceiverId string `json:"news-id" validate:"required"`
	Body       string `json:"body" validate:"max=150"`
}

type EventCommentRequest struct {
	AuthorId   string `json:"author-id" validate:"required"`
	ReceiverId string `json:"event-id" validate:"required"`
	Body       string `json:"body" validate:"max=150"`
}

type CreateEventRequest struct {
	AuthorId string `json:"author-id" validate:"required"`
	Game     string `json:"game" validate:"required"`
	Body     string `json:"body" validate:"max=150"`
	Max      int    `json:"max" validate:"required"`
	Minute   int    `json:"minute" validate:"required"`
}

type CreateNewsRequest struct {
	Title   string                `form:"title" validate:"required,max=100"`
	Body    string                `form:"body" validate:"required,max=150"`
	Link    string                `form:"link" validate:"required"`
	Picture *multipart.FileHeader `form:"picture" validate:"required"`
}

type PaginationRequest struct {
	Page   int `json:"page" validate:"required,gt=0"`
	Amount int `json:"amount" validate:"required,gt=0"`
}