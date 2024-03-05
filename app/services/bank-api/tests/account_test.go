package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	"github.com/qiushiyan/simplebank/business/core/account"
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
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, userAccount)
			},
		},
		{
			// try to get account that does not exist
			name:  "not-found",
			token: userToken,
			id:    1111,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(1111))).
					Times(1).
					Return(db_generated.Account{}, sql.ErrNoRows)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			// try to get account that does not belong to the user
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

		url := fmt.Sprintf("/accounts/%d", tc.id)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		request.Header.Add("authorization", tc.token)
		recorder := serveRequest(t, request, tc.buildStubs)
		tc.checker(recorder)
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
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
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

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, expected db_generated.Account) {
	got := getResponseData[db_generated.Account](t, body)
	require.Equal(t, got.ID, expected.ID)
	require.Equal(t, got.Owner, expected.Owner)
	require.Equal(t, got.Name, expected.Name)
	require.Equal(t, got.Balance, expected.Balance)
	require.Equal(t, got.Currency, expected.Currency)
	require.WithinDuration(t, got.CreatedAt, expected.CreatedAt, time.Second)
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
