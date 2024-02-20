package db

import (
	"context"
	"fmt"

	. "github.com/qiushiyan/simplebank/business/db/generated"
)

func (s *SQLStore) execTx(ctx context.Context, fn QueryFunc) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams is the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
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

func (s *SQLStore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}

		// create the entry for the from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		// create the entry for the to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// update two account balance
		var fromAccount, toAccount Account
		fromAccount, err = q.GetAccountForUpdate(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		toAccount, err = q.GetAccountForUpdate(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		// make sure we always update the lower account id first to avoid deadlocks
		// e.g., when two concurrent transaction go as a1 -> a2 and a2 -> a1
		// we always update a1 first before a2
		if fromAccount.ID < toAccount.ID {

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     fromAccount.ID,
				Amount: -args.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(
				ctx,
				AddAccountBalanceParams{ID: toAccount.ID, Amount: args.Amount},
			)
			if err != nil {
				return err
			}
		} else {

			result.ToAccount, err = q.AddAccountBalance(
				ctx,
				AddAccountBalanceParams{ID: toAccount.ID, Amount: args.Amount},
			)
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(
				ctx,
				AddAccountBalanceParams{ID: fromAccount.ID, Amount: -args.Amount},
			)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
