package tests

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	db "github.com/qiushiyan/simplebank/business/db/core"
	. "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
)

func createRandomAccount() Account {
	user, _ := createRandomUser()
	arg := CreateAccountParams{
		Owner:    user.Username,
		Name:     "test",
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
	user, _ := createRandomUser()
	arg := CreateAccountParams{
		Owner:    user.Username,
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

func TestUpdateAccountName(t *testing.T) {
	account1 := createRandomAccount()

	newName := "new name"
	arg := UpdateAccountNameParams{
		ID:   account1.ID,
		Name: "new name",
	}

	account2, err := testQueries.UpdateAccountName(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, newName, account2.Name)
}

func TestUpdateAccountBalance(t *testing.T) {
	account1 := createRandomAccount()

	amount := random.RandomMoney()
	arg := UpdateAccountBalanceParams{
		ID:     account1.ID,
		Amount: amount,
	}

	account2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.Balance+arg.Amount, account2.Balance)

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
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount()
	}

	owner := db.NewNullString(&lastAccount.Owner)
	arg := ListAccountsParams{
		Owner:  owner,
		Limit:  1,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.NotEmpty(t, accounts[0])
}
