package tests

import (
	"context"
	"testing"

	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/stretchr/testify/require"
)

func TestCreateGetFriend(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount()
	account2 := createRandomAccount()
	data, err := testQueries.CreateFriend(ctx, db_generated.CreateFriendParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, data)

	require.Equal(t, data.FromAccountID, account1.ID)
	require.Equal(t, data.ToAccountID, account2.ID)
	require.Equal(t, data.Status, "pending")

	record, err := testQueries.GetFriend(ctx, data.ID)
	require.NoError(t, err)
	require.Equal(t, record.ID, data.ID)
	require.Equal(t, record.FromAccountID, account1.ID)
	require.Equal(t, record.ToAccountID, account2.ID)
	require.Equal(t, record.Status, "pending")
}

func TestProcessFriend(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount()
	account2 := createRandomAccount()
	data, _ := testQueries.CreateFriend(ctx, db_generated.CreateFriendParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	})

	// accept
	res, err := testQueries.UpdateFriend(ctx, db_generated.UpdateFriendParams{
		ID:     data.ID,
		Status: "accepted",
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, "accepted")

	// reject
	res, err = testQueries.UpdateFriend(ctx, db_generated.UpdateFriendParams{
		ID:     data.ID,
		Status: "rejected",
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, "rejected")
}

func TestListFriends(t *testing.T) {
	fromAccount := createRandomAccount()
	ctx := context.Background()

	for i := range 10 {
		toAccount := createRandomAccount()
		data, _ := testQueries.CreateFriend(ctx, db_generated.CreateFriendParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
		})

		if i < 5 {
			testQueries.UpdateFriend(ctx, db_generated.UpdateFriendParams{
				ID:     data.ID,
				Status: "accepted",
			})
		}
	}

	// accepted requests
	status := "accepted"
	records, err := testQueries.ListFriends(ctx, db_generated.ListFriendsParams{
		FromAccountID: db.NewInt8(&fromAccount.ID),
		Status:        db.NewText(&status),
		Offset:        0,
		Limit:         5,
	})
	require.NoError(t, err)

	for i := range records {
		record := records[i]
		require.Equal(t, record.Status, "accepted")
	}
}
