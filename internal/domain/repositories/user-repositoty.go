package repositories

import (
	"context"
	"crap/internal/domain/entities"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface{
	Create(ctx context.Context, user entities.User) error
	FindById(ctx context.Context, id string) (*entities.User, error)
	Save(ctx context.Context, user entities.User) error
	ExistByLoginOrTg(ctx context.Context, login, tg string) (bool,error)
	Fetch(ctx context.Context, amount, page int) ([]entities.User, error)
	FindBy(ctx context.Context,vari, val string) (*entities.User, error)
}

type userRepository struct {
	DB *pgx.Conn
	Redis *redis.Client
}

func NewUserRepository(db *pgx.Conn, redis *redis.Client) UserRepository {
	return &userRepository{
		DB: db,
		Redis: redis,
	}
}

func (ur *userRepository) Create(ctx context.Context, user entities.User) error {
	if _,err := ur.DB.Exec(ctx,"INSERT INTO users (id,login,telegram,password,date_of_register) VALUES ($1,$2,$3,$4,$5)", user.Id,user.Login,user.Telegram,user.Password,user.DateOfRegister);err!=nil{
		return err
	}
	if ur.Redis != nil {
		userdata, err := json.Marshal(user)
		if err != nil {
			return err
		}
		ur.Redis.Set(ctx, user.Id.String(), userdata, time.Hour*24)
	}
	return nil
}

func (ur *userRepository) Save(ctx context.Context, user entities.User) error {
	if _,err := ur.DB.Exec(ctx,"UPDATE users SET chat_id=$1,rating=$2,total_rating=$3,number_of_rating=$4,games=$5,avatar=$6,discord=$7, date_of_register=$8 where id = $9",
	user.ChatId,user.Rating,user.TotalRating,user.NumberOfRating,user.Games,user.Avatar,user.Discord,user.DateOfRegister,user.Id);err!=nil {
		return err
	}
	if ur.Redis != nil {
		userdata, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := ur.Redis.Set(ctx, user.Id.String(), userdata, time.Hour*24).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (ur *userRepository) ExistByLoginOrTg(ctx context.Context, login, tg string) (bool,error) {
	var id string
	if err:=ur.DB.QueryRow(ctx,"SELECT id FROM users where login = $1 OR telegram = $2",login,tg).Scan(&id);err!=nil{
		if errors.Is(err, pgx.ErrNoRows){
			return false,nil
		}
		return false,err
	}
	return true,nil
}

func (ur *userRepository) FindBy(ctx context.Context,vari,val string) (*entities.User, error){
		user:=entities.User{}
		query:=fmt.Sprintf("SELECT id,login,telegram,chat_id,rating,total_rating,number_of_ratings,games,avatar,discord, date_of_register::timestamp from users where %s = $1",vari)
		if err := ur.DB.QueryRow(ctx,query,val).Scan(
			&user.Id,
			&user.Login,
			&user.Telegram,
			&user.ChatId,
			&user.Rating,
			&user.TotalRating,
			&user.NumberOfRating,
			&user.Games,
			&user.Avatar,
			&user.Discord,
			&user.DateOfRegister,
		)

		err != nil {
			return nil,err
		}

	return &user,nil
}

func (ur *userRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	user := &entities.User{}
	if ur.Redis != nil {
		userdata, err := ur.Redis.Get(ctx, id).Result()
		if err != nil {
			if err == redis.Nil {
					var err error
					user,err=ur.FindBy(ctx,"id",id)
					if err!=nil{
						return nil,err
					}
					userdata, err := json.Marshal(user)
					if err != nil {
						return nil, err
					}
					if err := ur.Redis.Set(ctx, id, userdata, time.Hour*24).Err(); err != nil {
						return nil, err
					}
				} else {
					var err error
					user,err=ur.FindBy(ctx,"id",id)
					if err!=nil{
						return nil,err
					}
				}
		} else {
			if err := json.Unmarshal([]byte(userdata), &user); err != nil {
				return nil, err
			}
		}
	} else {
		var err error
		user,err=ur.FindBy(ctx,"id",id)
		if err!=nil{
			return nil,err
		}
	}
	return user, nil
}

func (ur *userRepository) Fetch(ctx context.Context, amount, page int) ([]entities.User, error){
	users := []entities.User{}
	query := "SELECT * FROM users ORDER BY rating DESC OFFSET $1 LIMIT $2"
	rows, err := ur.DB.Query(ctx, query, page*amount-amount, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := entities.User{}
		if err := rows.Scan(
		&user.Id,
		&user.Login,
		&user.Telegram,
		&user.ChatId,
		&user.Rating,
		&user.TotalRating,
		&user.NumberOfRating,
		&user.Games,
		&user.Password,
		&user.Avatar,
		&user.Discord,
		)
		err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

