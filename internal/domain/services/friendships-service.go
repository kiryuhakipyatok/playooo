package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"crap/internal/dto"
)

type FriendshipsService interface {
	AddFriend(ctx context.Context, req dto.AddFriendRequest) error
	GetFriends(ctx context.Context, req dto.GetFriendsRequest) ([]entities.User, error)
	CancelFriendship(ctx context.Context, req dto.CancelFriendshipRequest) error
	AcceptFriendship(ctx context.Context, req dto.AcceptFriendshipRequest) error
	GetFriendRequests(ctx context.Context, req dto.GetFriendsReqRequests) ([]entities.User, error)
}

type friendshipsService struct{
	FriendshipsRepository repositories.FriendshipsRepository
	UserRepository repositories.UserRepository
}

func NewFriendshipsService(fr repositories.FriendshipsRepository, ur repositories.UserRepository) FriendshipsService{
	return &friendshipsService{
		FriendshipsRepository: fr,
		UserRepository: ur,
	}
}

func(fr *friendshipsService) AddFriend(ctx context.Context, req dto.AddFriendRequest) error{
	user1,err:=fr.UserRepository.FindById(ctx,req.UserId)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindBy(ctx,"login",req.FriendLogin)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Add(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) GetFriends(ctx context.Context, req dto.GetFriendsRequest) ([]entities.User, error){
	user,err:=fr.UserRepository.FindById(ctx,req.UserId)
	if err!=nil{
		return nil,err
	}
	friendsId,err:=fr.FriendshipsRepository.Fetch(ctx,user.Id.String(),req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	friends:=[]entities.User{}
	for _,fid:=range friendsId{
		friend,err:=fr.UserRepository.FindById(ctx,fid)
		if err!=nil{
			return nil,err
		}
		friends=append(friends, *friend)
	}
	return friends,nil
}

func(fr *friendshipsService) CancelFriendship(ctx context.Context, req dto.CancelFriendshipRequest) error{
	user1,err:=fr.UserRepository.FindById(ctx,req.UserId)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindById(ctx,req.FriendId)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Cancel(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) AcceptFriendship(ctx context.Context, req dto.AcceptFriendshipRequest) error{
	user1,err:=fr.UserRepository.FindById(ctx,req.UserId)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindById(ctx, req.FriendId)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Accept(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) GetFriendRequests(ctx context.Context, req dto.GetFriendsReqRequests) ([]entities.User, error){
	user,err:=fr.UserRepository.FindById(ctx,req.UserId)
	if err!=nil{
		return nil,err
	}
	requestsId,err:=fr.FriendshipsRepository.FetchRequests(ctx,user.Id.String(),req.Amount,req.Page)
	if err!=nil{
		return nil,err
	}
	requests:=[]entities.User{}
	for _,rid:=range requestsId{
		request,err:=fr.UserRepository.FindById(ctx,rid)
		if err!=nil{
			return nil,err
		}
		requests=append(requests, *request)
	}
	return requests,nil
}

