package db

import (
	"context"

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
