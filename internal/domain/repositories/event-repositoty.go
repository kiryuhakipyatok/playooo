package repositories

import (
	"context"
	"crap/internal/domain/entities"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type EventRepository interface{
	Create(ctx context.Context, event entities.Event) error
	Delete(ctx context.Context, event entities.Event) error
	FindById(ctx context.Context, id string) (*entities.Event, error)
	FetchUpcoming(ctx context.Context, time time.Time) ([]entities.Event, error)
	Fetch(ctx context.Context, amount, page int) ([]entities.Event, error)
	Join(ctx context.Context,user_id,event_id string) error 
	Unjoin(ctx context.Context,user_id,event_id string) error
	FetchMembers(ctx context.Context,id string) ([]string,error)
	Save(c context.Context, event entities.Event) error
	Filter(ctx context.Context, game,max,time string, amount, page int) ([]entities.Event, error)
	Sort(ctx context.Context, field,dir string, amount, page int) ([]entities.Event, error)
}

type eventRepository struct {
	DB    *pgx.Conn
	Redis *redis.Client
}

func NewEventRepository(db *pgx.Conn, redis *redis.Client) EventRepository {
	return &eventRepository{
		DB:    db,
		Redis: redis,
	}
}

func (er *eventRepository) Create(ctx context.Context, event entities.Event) error {
	if _,err := er.DB.Exec(ctx, "INSERT INTO events (id,author_id,body,game,max,time,notificated_pre) values($1,$2,$3,$4,$5,$6,$7)", event.Id, event.AuthorId, event.Body, event.Game, event.Max, event.Time, event.NotificatedPre); err != nil {
		return err
	}
	if _,err:=er.DB.Exec(ctx,"INSERT INTO users_events (event_id,user_id) values($1,$2)",event.Id,event.AuthorId);err!=nil{
		return err
	}
	eventdata, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if er.Redis != nil {
		if err := er.Redis.Set(ctx, event.Id.String(), eventdata, time.Until(event.Time)).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (er *eventRepository) Save(ctx context.Context, event entities.Event) error {
	if _,err := er.DB.Exec(ctx, "UPDATE events SET author_id=$1,body=$2,game=$3,max=$4,time=$5,notificated_pre=$6 WHERE id = $7",event.AuthorId,event.Body, event.Game, event.Max, event.Time, event.NotificatedPre,event.Id); err != nil {
		return err
	}
	if er.Redis != nil {
		eventdata, err := json.Marshal(event)
		if err != nil {
			return err
		}
		ttl, err := er.Redis.TTL(ctx, event.Id.String()).Result()
		if err != nil {
			return err
		}
		if err := er.Redis.Set(ctx, event.Id.String(), eventdata, ttl).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (er *eventRepository) Delete(ctx context.Context, event entities.Event) error {
	if _,err:=er.DB.Exec(ctx,"DELETE FROM events WHERE id = $1",event.Id);err!=nil{
		return err
	}
	if _,err:=er.DB.Exec(ctx,"DELETE FROM users_events WHERE event_id = $1",event.Id);err!=nil{
		return err
	}
	return nil
}

func (er *eventRepository) FetchUpcoming(ctx context.Context, time time.Time) ([]entities.Event, error) {
	events := []entities.Event{}
	rows,err:=er.DB.Query(ctx,"SELECT * FROM events where time <= $1",time)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		event:=entities.Event{}
		if err:=rows.Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
			return nil,err
		}
		events=append(events, event)
	}
	return events, nil
}

func (er *eventRepository) FindById(ctx context.Context, id string) (*entities.Event, error) {
	event := entities.Event{}
	if er.Redis != nil {
		eventdata, err := er.Redis.Get(ctx, id).Result()
		if err != nil {
			if err == redis.Nil {
				if err:=er.DB.QueryRow(ctx,"SELECT * FROM events where id= $1",id).Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
					return nil,err
				}
				eventdata, err := json.Marshal(event)
				if err != nil {
					return nil, err
				}
				if err := er.Redis.Set(ctx, id, eventdata, time.Until(event.Time)).Err(); err != nil {
					return nil, err
				}
			} else {
				if err:=er.DB.QueryRow(ctx,"SELECT * FROM events where id= $1",id).Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
					return nil,err
				}
			}
		} else {
			if err := json.Unmarshal([]byte(eventdata), &event); err != nil {
				return nil, err
			}
		}
	} else {
		if err:=er.DB.QueryRow(ctx,"SELECT * FROM events where id= $1",id).Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
			return nil,err
		}
	}
	return &event, nil
}

func (er *eventRepository) Fetch(ctx context.Context, amount, page int) ([]entities.Event, error) {
	events := []entities.Event{}
	query := "SELECT * FROM events OFFSET $1 LIMIT $2"
	rows, err := er.DB.Query(ctx, query, page*amount-amount, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		event:=entities.Event{}
		if err:=rows.Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
			return nil,err
		}
		events=append(events, event)
	}
	return events, nil
}

func (er *eventRepository) Join(ctx context.Context,user_id,event_id string) error{
	if _,err:=er.DB.Exec(ctx,"INSERT INTO users_events (event_id,user_id) values($1,$2)",event_id,user_id);err!=nil{
		return err
	}
	return nil
}

func (er *eventRepository) Unjoin(ctx context.Context,user_id,event_id string) error{
	if _,err:=er.DB.Exec(ctx,"DELETE FROM users_events where event_id=$1 AND user_id = $2 ",event_id,user_id);err!=nil{
		return err
	}
	return nil
}

func (er *eventRepository) FetchMembers(ctx context.Context,id string) ([]string,error){
	members:=[]string{}
	rows,err:=er.DB.Query(ctx,"SELECT user_id FROM users_events WHERE event_id = $1",id)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		var id string
		if err:=rows.Scan(&id);err!=nil{
			return nil,err
		}
		members = append(members, id)
	}
	return members,nil
}

func (er *eventRepository) Filter(ctx context.Context, game,max,time string, amount, page int) ([]entities.Event, error){
	var q string
	fields:=map[string]string{
		"game":game,
		"max":max,
		"time":time,
	}
	for f,v:=range fields{
		if v!=""{
			if len(q)!=0{
				q+=fmt.Sprintf(" AND %s='%s'",f,v)
			}else{
				q+=fmt.Sprintf(" WHERE %s='%s'",f,v)
			}
		}
	}
	events:=[]entities.Event{}
	query:=fmt.Sprintf("SELECT * FROM events %s OFFSET $2 LIMIT $3",q)
	rows,err:=er.DB.Query(ctx,query,amount*page-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		event:=entities.Event{}
		if err:=rows.Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
			return nil,err
		}
		events=append(events,event)
	}
	return events,nil
}

func (er *eventRepository) Sort(ctx context.Context, field,dir string, amount, page int) ([]entities.Event, error){
	events:=[]entities.Event{}
	query:=fmt.Sprintf("SELECT * FROM events ORDER BY %s %s OFFSET $1 LIMIT $2",field,dir)
	rows,err:=er.DB.Query(ctx, query,amount*page-amount,amount)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		event:=entities.Event{}
		if err:=rows.Scan(&event.Id,&event.AuthorId,&event.Body,&event.Game,&event.Max,&event.Time,&event.NotificatedPre);err!=nil{
			return nil,err
		}
		events=append(events,event)
	}
	return events,nil
}