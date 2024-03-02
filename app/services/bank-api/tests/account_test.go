package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/foundation/validate"
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
			id:    userAccount.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(userAccount.ID)).
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
			id:    adminAccount.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(adminAccount.ID)).
					Times(1).
					Return(adminAccount, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
	newAccount := db_generated.Account{
		Owner:    "user",
		Currency: validate.CAD,
		Balance:  0,
	}
	args := db_generated.CreateAccountParams{
		Owner:    newAccount.Owner,
		Currency: newAccount.Currency,
		Balance:  newAccount.Balance,
	}
	body := accountgrp.CreateAccountRequest{
		Currency: newAccount.Currency,
	}
	cases := []struct {
		name       string
		currency   string
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			// create new account as user
			name:     "ok",
			token:    userToken,
			currency: newAccount.Currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(newAccount, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, newAccount)
			},
		},
		{
			// create new account without token
			name:     "unauthorized",
			token:    "",
			currency: newAccount.Currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(args)).
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
	require.Equal(t, got, expected)
}

var userAccount = db_generated.Account{
	ID:       3,
	Owner:    "user",
	Currency: "USD",
	Balance:  100,
}

var adminAccount = db_generated.Account{
	ID:       1,
	Owner:    "admin",
	Currency: "USD",
	Balance:  100,
}
