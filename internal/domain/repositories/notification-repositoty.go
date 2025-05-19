package repositories

import (
	"context"
	"crap/internal/domain/entities"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification entities.Notification) error
	CreateForUsers(ctx context.Context, notification entities.Notification, id string ) error
	Delete(ctx context.Context, id string, nid string) error
	DeleteAll(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*entities.Notification, error)
	Fetch(ctx context.Context,id string, amount, page int) ([]entities.Notification, error)
}

type notificationRepository struct {
	DB  *pgx.Conn
}

func NewNoticeRepository(db *pgx.Conn) NotificationRepository {
	return &notificationRepository{
		DB: db,
	}
}

func (nr *notificationRepository) Create(ctx context.Context, notification entities.Notification) error {
	fmt.Println(notification)
	if _,err:=nr.DB.Exec(ctx,"INSERT INTO notifications (id,event_id,body,time) values($1,$2,$3,$4)",notification.Id,notification.EventId,notification.Body,notification.Time);err!=nil{
		return err
	}
	
	return nil
}

func (nr *notificationRepository) CreateForUsers(ctx context.Context, notification entities.Notification, id string ) error{
	if _,err:=nr.DB.Exec(ctx,"INSERT INTO users_notifications (user_id,notification_id) values($1,$2)",id,notification.Id);err!=nil{
		return err
	}
	return nil
}

func (nr *notificationRepository) Delete(ctx context.Context, id, nid string) error {
	if _,err:=nr.DB.Exec(ctx,"DELETE FROM users_notifications WHERE user_id = $1 AND notification_id=$2",id,nid);err!=nil{
		return err
	}
	return nil
}

func (nr *notificationRepository) DeleteAll(ctx context.Context, id string) error{
	if _,err:=nr.DB.Exec(ctx,"DELETE FROM users_notifications WHERE user_id = $1)",id);err!=nil{
		return err
	}
	return nil
}

func (nr *notificationRepository) FindById(ctx context.Context, id string) (*entities.Notification, error) {
	notification := entities.Notification{}
	if err:=nr.DB.QueryRow(ctx,"SELECT * FROM notifications WHERE id = $1",id).Scan(&notification.Id,&notification.EventId,&notification.Body,&notification.Time);err!=nil{
		return nil,err
	}
	return &notification, nil
}

func (nr *notificationRepository) Fetch(ctx context.Context,id string, amount, page int) ([]entities.Notification, error){
	notifications:=[]entities.Notification{}
	query:="SELECT * FROM notifications n JOIN users_notifications un ON un.notification_id = n.id WHERE un.user_id = $1 ORDER BY n.time OFFSET $2 LIMIT $3"
	rows,err:=nr.DB.Query(ctx,query,id,page*amount-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		n:=entities.Notification{}
		if err:=rows.Scan(&n.Id,&n.EventId,&n.Body,&n.Time);err!=nil{
			return nil,err
		}
		notifications=append(notifications,n)
	}
	if rows.Err()!=nil{
		return nil,err
	}
	return notifications,nil
}