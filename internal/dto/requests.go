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

type CreateEventRequest struct {
	AuthorId string `json:"author-id" validate:"required"`
	Game     string `json:"game" validate:"required"`
	Body     string `json:"body" validate:"max=150"`
	Max      int    `json:"max" validate:"required"`
	Minute   int    `json:"minute" validate:"required"`
}

type JoinToEventRequest struct{
	UserId string `json:"user-id" validate:"required"`
	EventId string `json:"event-id" validate:"required"`
}

type UnjoinFromEventRequest struct{
	JoinToEventRequest
}

type AddFriendRequest struct{
	UserId string `json:"user-id" validate:"required"`
	FriendLogin string `json:"friend-login" validate:"required"`
}

type AcceptFriendshipRequest struct{
	UserId string `json:"user-id" validate:"required"`
	FriendId string `json:"friend-id" validate:"required"`
}

type CancelFriendshipRequest struct{
	AcceptFriendshipRequest
}

type GetFriendsRequest struct{
	UserId string `json:"user-id" validate:"required"`
	PaginationRequest
}

type GetNotificationsRequest struct{
	UserId string `json:"user-id" validate:"required"`
	PaginationRequest
}

type GetFriendsReqRequests struct{
	GetFriendsRequest
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

type AddCommentRequest struct{
	Whom string `json:"whom" validate:"required,max=6"`
	UserId string `json:"user-id" validate:"required"`
	ReceiverId string `json:"receiver-id" validate:"required"`
	Body string `json:"body" validate:"required,max=150"`
}

type GetCommentsRequest struct{
	Whose string `json:"whom" validate:"required,max=6"`
	UserId string `json:"user-id" validate:"required"`
	PaginationRequest
}

type AddGameRequest struct{
	UserId string `json:"user-id" validate:"required"`
	Game string `json:"game" validate:"required"`
}

type DeleteGameRequest struct{
	AddGameRequest
}

type UploadAvatarRequest struct{
	UserId string `json:"user-id" validate:"required"`
	Picture *multipart.FileHeader `json:"picture" validate:"required"`
}

type RecordDiscordRequest struct{
	UserId string `json:"user-id" validate:"required"`
	Discord string `json:"discord" validate:"required"`
}

type EditRatingRequest struct{
	UserId string `json:"user-id" validate:"required"`
	Stars int `json:"stars" validate:"required"`
}

