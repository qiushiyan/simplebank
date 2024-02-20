package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	. "github.com/qiushiyan/simplebank/business/db/generated"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type QueryFunc = func(*Queries) error

// SQLStore is the implementation of the Store interface that uses postgres
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func Open(user, password, host, port, dbname string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	return sql.Open("postgres", connectionString)
}
