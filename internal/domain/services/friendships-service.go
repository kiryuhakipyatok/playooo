package services

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
)

type FriendshipsService interface {
	AddFriend(ctx context.Context, id, login string) error
	ShowFriends(ctx context.Context, id string, amount, page int) ([]entities.User, error)
	CancelFriendship(ctx context.Context, id, login string) error
	AcceptFriendship(ctx context.Context, id, login string) error
	GetFriendRequests(ctx context.Context, id string, amount, page int) ([]entities.User, error)
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

func(fr *friendshipsService) AddFriend(ctx context.Context, id, login string) error{
	user1,err:=fr.UserRepository.FindById(ctx,id)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindBy(ctx,"login",login)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Add(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) ShowFriends(ctx context.Context, id string, amount, page int) ([]entities.User, error){
	user,err:=fr.UserRepository.FindById(ctx,id)
	if err!=nil{
		return nil,err
	}
	friendsId,err:=fr.FriendshipsRepository.Fetch(ctx,user.Id.String(),amount,page)
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

func(fr *friendshipsService) CancelFriendship(ctx context.Context, id, login string) error{
	user1,err:=fr.UserRepository.FindById(ctx,id)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindBy(ctx,"login",login)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Cancel(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) AcceptFriendship(ctx context.Context, id, login string) error{
	user1,err:=fr.UserRepository.FindById(ctx,id)
	if err!=nil{
		return err
	}
	user2,err:=fr.UserRepository.FindBy(ctx,"login",login)
	if err!=nil{
		return err
	}
	if err:=fr.FriendshipsRepository.Accept(ctx,user1.Id.String(),user2.Id.String());err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsService) GetFriendRequests(ctx context.Context, id string, amount, page int) ([]entities.User, error){
	user,err:=fr.UserRepository.FindById(ctx,id)
	if err!=nil{
		return nil,err
	}
	requestsId,err:=fr.FriendshipsRepository.FetchRequests(ctx,user.Id.String(),amount,page)
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

