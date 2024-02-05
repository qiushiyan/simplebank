package db

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	. "github.com/qiushiyan/bank-api/business/db/generated"
	"github.com/qiushiyan/bank-api/business/random"
	"github.com/stretchr/testify/require"
)

func createRandomAccount() Account {
	arg := CreateAccountParams{
		Owner:    random.RandomOwner(),
		Currency: random.RandomCurrency(),
		Balance:  random.RandomMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		log.Fatalf("cannot create account: %v", err)
	}
	return account
}

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    random.RandomOwner(),
		Currency: random.RandomCurrency(),
		Balance:  random.RandomMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount()
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount()

	b := random.RandomMoney()
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: b,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, b, account2.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount()
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount()
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for i := range accounts {
		a := accounts[i]
		require.NotEmpty(t, a)
	}

}
