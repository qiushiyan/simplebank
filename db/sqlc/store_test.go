package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	var wg sync.WaitGroup
	store := NewStore(testDB)

	a1 := createRandomAccount()
	a2 := createRandomAccount()

	n := 5
	amount := int64(10)
	ctx := context.Background()

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := store.TransferTx(ctx, TransferTxParams{
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

		fmt.Println("k is", k)

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
