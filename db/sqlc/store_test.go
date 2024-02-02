package db

import (
	"context"
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

	for r := range results {
		require.NotEmpty(t, r)

		// check transfer
		require.Equal(t, a1.ID, r.Transfer.FromAccountID)
		require.Equal(t, a2.ID, r.Transfer.ToAccountID)
		require.Equal(t, amount, r.Transfer.Amount)

		_, err := testQueries.GetTransfer(ctx, r.Transfer.ID)
		require.NoError(t, err)
		require.Equal(t, a1.ID, r.Transfer.FromAccountID)
		require.Equal(t, a2.ID, r.Transfer.ToAccountID)

		// check entries
		fromEntry, err := testQueries.GetEntry(ctx, r.FromEntry.ID)
		require.NoError(t, err)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		toEntry, err := testQueries.GetEntry(ctx, r.ToEntry.ID)
		require.NoError(t, err)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// TODO: check account balances
	}

	// Check the initial account balances
	// Transfer some money between the accounts
	// Check the new account balances
	// Check the transfer record in the database
}
