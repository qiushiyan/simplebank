package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	. "github.com/qiushiyan/simplebank/business/db/generated"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Check(ctx context.Context) error
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

// Checl returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func (s *SQLStore) Check(ctx context.Context) error {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = s.db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return s.db.QueryRowContext(ctx, q).Scan(&tmp)
}

func Open(user, password, host, port, dbname string) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)

	return sql.Open("postgres", connectionString)
}
