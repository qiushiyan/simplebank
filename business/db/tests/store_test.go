package tests

import (
	"context"
	"sync"
	"testing"

	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/stretchr/testify/require"
)

// check for transfer (single direction)
func TestTransferTx(t *testing.T) {
	var wg sync.WaitGroup
	store := db.NewPostgresStore(testDB)

	a1 := createRandomAccount()
	a2 := createRandomAccount()

	n := 5
	amount := int64(10)
	ctx := context.Background()

	errs := make(chan error, n)
	results := make(chan db.TransferTxResult, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := store.TransferTx(ctx, db.TransferTxParams{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	wg.Wait()
	close(errs)
	close(results)

	for err := range errs {
		require.NoError(t, err)
	}

	var existed = make(map[int]bool)
	for r := range results {
		require.NotEmpty(t, r)

		// check transfer
		transfer := r.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, a1.ID, transfer.FromAccountID)
		require.Equal(t, a2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, a1.ID, r.Transfer.FromAccountID)
		require.Equal(t, a2.ID, r.Transfer.ToAccountID)

		// check entries
		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// check accounts
		fromAccount := r.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, a1.ID, fromAccount.ID)

		toAccount := r.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, a2.ID, toAccount.ID)

		// check balances
		diff1 := a1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - a2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
		// check k is unique across transactions and can be only picked from 1, ..., n
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedFromAccount, err := testQueries.GetAccount(ctx, a1.ID)
	require.NoError(t, err)
	require.Equal(t, a1.Balance-int64(n)*amount, updatedFromAccount.Balance)

	updatedToAccount, err := testQueries.GetAccount(ctx, a2.ID)
	require.NoError(t, err)
	require.Equal(t, a2.Balance+int64(n)*amount, updatedToAccount.Balance)
}

// simulate two accounts transferring to each other concurrently
func TestTransferTxDeadlock(t *testing.T) {
	store := db.NewPostgresStore(testDB)

	a1 := createRandomAccount()
	a2 := createRandomAccount()

	n := 4
	amount := int64(10)
	ctx := context.Background()

	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func(i int) {
			var err error

			if i%2 == 1 {
				// transfer from a1 to a2
				_, err = store.TransferTx(ctx, db.TransferTxParams{
					FromAccountID: a1.ID,
					ToAccountID:   a2.ID,
					Amount:        amount,
				})

			} else {
				// transfer from a2 to a1

				_, err = store.TransferTx(ctx, db.TransferTxParams{
					FromAccountID: a2.ID,
					ToAccountID:   a1.ID,
					Amount:        amount,
				})
			}

			errs <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		require.NoError(t, <-errs)
	}

	// check the final balance, should be the same
	updatedFromAccount, err := testQueries.GetAccount(ctx, a1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(ctx, a2.ID)
	require.NoError(t, err)

	require.Equal(t, a1.Balance, updatedFromAccount.Balance)
	require.Equal(t, a2.Balance, updatedToAccount.Balance)
}
