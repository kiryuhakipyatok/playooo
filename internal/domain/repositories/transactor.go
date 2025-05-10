package repositories

import (
	"context"
	"log"
	"time"
	"github.com/jackc/pgx/v5"
)

type Transactor interface {
	WithinTransaction(context.Context, func(c context.Context) (any, error)) (any, error)
}

type transactor struct {
	DB *pgx.Conn
}

func NewTransactor(db *pgx.Conn) Transactor {
	return &transactor{
		DB: db,
	}
}

type txKey struct{}

func injectTx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func (t *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) (any, error)) (any, error) {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx,err := t.DB.BeginTx(txCtx, pgx.TxOptions{})
	if err!=nil{
		return nil,err
	}
	ctxTime,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	res, err := tFunc(injectTx(ctx, &tx))
	if err != nil {
		if err := tx.Rollback(ctxTime); err != nil {
			log.Printf("cannot rollback transaction: %v", err)
		}
		return nil, err
	}

	if err := tx.Commit(ctxTime); err != nil {
		log.Printf("cannot commit transaction: %v", err)
	}

	return res, nil
}
