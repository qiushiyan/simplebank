package db

import (
	"context"

	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	db_generated.Querier
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Check(ctx context.Context) error
}
