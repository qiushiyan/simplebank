package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/qiushiyan/simplebank/business/db/generated"
)

// PostgresStore is the implementation of the Store interface that uses the postgres database
type PostgresStore struct {
	*Queries
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		Queries: New(pool),
		pool:    pool,
	}
}

// Close closes all connections in the pool
func (s *PostgresStore) Close() {
	s.pool.Close()
}

// Check returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func (s *PostgresStore) Check(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {

		pingError = s.pool.Ping(ctx)
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
	return s.pool.QueryRow(ctx, q).Scan(&tmp)
}

func NewPgxPool(ctx context.Context, config string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.NewWithConfig(ctx, conf)
	return db, err
}

type AfterCreateUserFunc func(user User) (any, error)

// CreateUserTxParams is the input parameters for the create user function
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate AfterCreateUserFunc
}

// CreateUserTxResult is the result of the create user function
type CreateUserTxResult struct {
	User              User `json:"user"`
	AfterCreateResult any  `json:"result"`
}

func (s *PostgresStore) CreateUserTx(
	ctx context.Context,
	arg CreateUserTxParams,
) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.ExecuteInTransaction(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		if arg.AfterCreate == nil {
			result.AfterCreateResult = nil
			return nil
		}

		result.AfterCreateResult, err = arg.AfterCreate(result.User)
		return err
	})

	return result, err
}

// TransferTxParams is the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, create an entry for both accounts, and update accounts' balance within a single database transaction
func (s *PostgresStore) TransferTx(
	ctx context.Context,
	arg TransferTxParams,
) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.ExecuteInTransaction(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				arg.FromAccountId,
				-arg.Amount,
				arg.ToAccountId,
				arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

func (s *PostgresStore) ExecuteInTransaction(ctx context.Context, fn QueryFunc) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
