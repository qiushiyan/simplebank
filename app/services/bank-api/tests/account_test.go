package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	"github.com/qiushiyan/simplebank/business/core/account"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountApi(t *testing.T) {
	cases := []struct {
		name       string
		id         int64
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: userToken,
			id:    userAccountId,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(userAccountId)).
					Times(1).
					Return(userAccount, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				got := getResponseData[db_generated.Account](t, recorder.Body)
				requireMatchAccount(t, got, userAccount)
			},
		},
		{
			// try to get account that does not exist
			name:  "not-found",
			token: userToken,
			id:    0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(0))).
					Times(1).
					Return(db_generated.Account{}, pgx.ErrNoRows)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			// try to get account of a different user
			name:  "wrong-owner",
			token: userToken,
			id:    adminAccountId,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(adminAccountId)).
					Times(1).
					Return(adminAccount, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			url := fmt.Sprintf("/accounts/%d", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Add("authorization", tc.token)
			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})

	}
}

func TestCreateAccountApi(t *testing.T) {
	newAccount := account.NewAccount{
		Owner:    "user",
		Name:     "new account",
		Balance:  0,
		Currency: "USD",
	}
	arg := db_generated.CreateAccountParams(newAccount)

	account := db_generated.Account{
		Owner:     newAccount.Owner,
		Name:      newAccount.Name,
		Balance:   newAccount.Balance,
		Currency:  newAccount.Currency,
		CreatedAt: time.Now(),
	}

	body := accountgrp.CreateAccountRequest{
		Name:     newAccount.Name,
		Currency: newAccount.Currency,
	}
	cases := []struct {
		name       string
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			// create new account as user
			name:  "ok",
			token: userToken,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				got := getResponseData[db_generated.Account](t, recorder.Body)
				requireMatchAccount(t, got, account)
			},
		},
		{
			// create new account without token
			name:  "unauthorized",
			token: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b, err := json.Marshal(body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(b))
			require.NoError(t, err)

			request.Header.Add("authorization", tc.token)
			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})

	}
}

func TestListAccountsApi(t *testing.T) {
	newArg := func(owner string) db_generated.ListAccountsParams {
		return db_generated.ListAccountsParams{
			Limit:  5,
			Offset: 0,
			Owner:  db.NewText(&owner),
		}
	}

	cases := []struct {
		name       string
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			// list accounts as user
			name:  "ok",
			token: userToken,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), newArg(userAccount.Owner)).
					Times(1).
					Return([]db_generated.Account{userAccount}, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				got := getResponseData[[]db_generated.Account](t, recorder.Body)
				requireMatchAccounts(t, got, []db_generated.Account{userAccount})
			},
		},
		{
			// list accounts without token
			name:  "unauthorized",
			token: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			require.NoError(t, err)

			request.Header.Add("authorization", tc.token)
			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})
	}

}

func requireMatchAccount(
	t *testing.T,
	got db_generated.Account,
	expected db_generated.Account,
) {

	require.Equal(t, got.ID, expected.ID)
	require.Equal(t, got.Owner, expected.Owner)
	require.Equal(t, got.Name, expected.Name)
	require.Equal(t, got.Balance, expected.Balance)
	require.Equal(t, got.Currency, expected.Currency)
	require.WithinDuration(t, got.CreatedAt, expected.CreatedAt, time.Second)
}

func requireMatchAccounts(
	t *testing.T,
	got []db_generated.Account,
	expected []db_generated.Account,
) {
	require.True(t, len(got) == len(expected))
	for i := range got {
		requireMatchAccount(t, got[i], expected[i])
	}
}

var userAccountId int64 = 3
var userAccount = db_generated.Account{
	ID:       userAccountId,
	Owner:    "user",
	Currency: "USD",
	Balance:  100,
}

var adminAccountId int64 = 1
var adminAccount = db_generated.Account{
	ID:       adminAccountId,
	Owner:    "admin",
	Currency: "USD",
	Balance:  100,
}
