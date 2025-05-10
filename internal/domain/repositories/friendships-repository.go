package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"errors"
)

type FriendshipsRepository interface {
	Add(ctx context.Context,id1,id2 string) error
	Cancel(ctx context.Context, id1, id2 string) error
	Accept(ctx context.Context, id1, id2 string) error
	Fetch(ctx context.Context, id string, amount, page int) ([]string,error)
	FetchRequests(ctx context.Context, id string, amount, page int) ([]string,error)
}

type friendshipsRepository struct{
	DB *pgx.Conn
}

func NewFriendshipsRepository(db *pgx.Conn) FriendshipsRepository{
	return &friendshipsRepository{
		DB: db,
	}
}

func(fr *friendshipsRepository)	Add(ctx context.Context,id1,id2 string) error{
	var relation string
	if err:=fr.DB.QueryRow(ctx,"SELECT relation FROM friendships where (user_id1 = $1 and user_id2 = $2) OR (user_id2 = $1 and user_id1 = $2)",id1,id2).Scan(&relation);err!=nil{
		return err
	}
	switch relation{
	case "accepted":
		return errors.New("users already are friends")
	case "requested":
		return errors.New("friend request already sent")
	default:
		if _,err:=fr.DB.Exec(ctx,"INSERT INTO friendships (user_id1,user_id2) values($1,$2)",id1,id2);err!=nil{
			return err
		}
	}
	return nil
}

func(fr *friendshipsRepository)	Cancel(ctx context.Context, id1, id2 string) error{
	if _,err:=fr.DB.Exec(ctx,"DELETE FROM friendships WHERE (user_id1 = $1 and user_id2 = $2) OR (user_id1=$2 and user_id2=$1)",id1,id2);err!=nil{
		return err
	}
	return nil
}

func(fr *friendshipsRepository)	Accept(ctx context.Context, id1, id2 string) error{
	var relation string
	if err:=fr.DB.QueryRow(ctx,"SELECT relation FROM friendships where (user_id1 = $2 and user_id2 = $1)",id1,id2).Scan(&relation);err!=nil{
		return err
	}
	switch relation{
	case "accepted":
		return errors.New("users already are friends")
	case "":
		return errors.New("no friend request")
	default:
		if _,err:=fr.DB.Exec(ctx,"UPDATE friendships SET relation = 'accepted' where (user_id2 = $1 and user_id1 = $2)",id1,id2);err!=nil{
			return err
	 	}
	}
	return nil
}

func(fr *friendshipsRepository)	Fetch(ctx context.Context, id string, amount, page int) ([]string,error){
	friends:=[]string{}
	rows,err:=fr.DB.Query(ctx,"SELECT * FROM users u JOIN friendships f ON f.user_id2=u.id WHERE f.user_id1=$1 AND f.relation = 'accepted' UNION SELECT u.id FROM users u JOIN friendships f ON f.user_id1=u.id WHERE f.user_id2=$1 AND f.relation = 'accepted' OFFSET $2 LIMIT $3",id,page*amount-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		var id string
		if err:=rows.Scan(&id);err!=nil{
			return nil,err
		}
		friends=append(friends, id)
	}
	return friends,nil
}

func(fr *friendshipsRepository)	FetchRequests(ctx context.Context, id string, amount, page int) ([]string,error){
	requests:=[]string{}
	rows,err:=fr.DB.Query(ctx,"SELECT user_id1 FROM friendships WHERE user_id2 = $1 AND relation='requested'",id)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		var id string
		if err:=rows.Scan(&id);err!=nil{
			return nil,err
		}
		requests = append(requests, id)
	}
	return requests,nil
}
