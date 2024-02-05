package db

import (
	"context"
	"testing"
	"time"

	. "github.com/qiushiyan/bank-api/business/db/generated"
	"github.com/qiushiyan/bank-api/business/random"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount()

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount()

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, entry.Amount, entry2.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, 1*time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount()

	createArgs := CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.RandomMoney(),
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}

	for i := 0; i < 10; i++ {
		testQueries.CreateEntry(context.Background(), createArgs)
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	if err != nil {
		t.Error(err)
	}

	for i := range entries {
		require.NotEmpty(t, entries[i])
		require.Equal(t, createArgs.AccountID, entries[i].AccountID)
		require.Equal(t, createArgs.Amount, entries[i].Amount)
	}
}
