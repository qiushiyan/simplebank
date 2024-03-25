package tests

import (
	"context"
	"testing"

	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/stretchr/testify/require"
)

func TestCreateFriend(t *testing.T) {
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
	require.True(t, data.Pending)
	require.False(t, data.Accepted)
}

func TestProcessFriend(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount()
	account2 := createRandomAccount()
	data, _ := testQueries.CreateFriend(ctx, db_generated.CreateFriendParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	})

	res, err := testQueries.AcceptFriend(ctx, data.ID)
	require.NoError(t, err)
	require.False(t, res.Pending)
	require.True(t, res.Accepted)

	res, err = testQueries.DeclineFriend(ctx, data.ID)
	require.NoError(t, err)
	require.False(t, res.Accepted)
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
			testQueries.AcceptFriend(ctx, data.ID)
		}
	}

	// accepted requests
	accepted := true
	records, err := testQueries.ListFriends(ctx, db_generated.ListFriendsParams{
		FromAccountID: db.NewInt8(&fromAccount.ID),
		Pending:       db.NewBool(nil),
		Accepted:      db.NewBool(&accepted),
		Offset:        0,
		Limit:         5,
	})
	require.NoError(t, err)

	for i := range records {
		record := records[i]
		require.True(t, record.Accepted)
		require.False(t, record.Pending)
	}

}
