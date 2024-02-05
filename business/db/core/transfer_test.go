package db

import (
	"context"
	"testing"

	. "github.com/qiushiyan/bank-api/business/db/generated"
	"github.com/qiushiyan/bank-api/business/random"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	a1 := createRandomAccount()
	a2 := createRandomAccount()

	args := CreateTransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        random.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
}

func TestGetTransfer(t *testing.T) {
	a1 := createRandomAccount()
	a2 := createRandomAccount()

	args := CreateTransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        random.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	if err != nil {
		t.Error(err)
	}

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.Amount, transfer2.Amount)
}

func TestListTransfers(t *testing.T) {
	a1 := createRandomAccount()
	a2 := createRandomAccount()

	createArgs := CreateTransferParams{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        random.RandomMoney(),
	}

	args := ListTransfersParams{
		FromAccountID: a1.ID,
		Limit:         5,
		Offset:        5,
	}

	for i := 0; i < 10; i++ {
		testQueries.CreateTransfer(context.Background(), createArgs)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	if err != nil {
		t.Error(err)
	}

	for i := range transfers {
		require.NotEmpty(t, transfers[i])
		require.Equal(t, createArgs.FromAccountID, transfers[i].FromAccountID)
		require.Equal(t, createArgs.ToAccountID, transfers[i].ToAccountID)
		require.Equal(t, createArgs.Amount, transfers[i].Amount)
	}
}
