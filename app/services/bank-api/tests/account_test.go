package tests

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	db "github.com/qiushiyan/simplebank/business/db/core"
	. "github.com/qiushiyan/simplebank/business/db/generated"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	cases := []struct {
		name       string
		url        string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder, error)
	}{
		{
			name: "ok",
			url:  fmt.Sprintf("/accounts/%d", account.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder, err error) {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "not-found",
			url:  "/accounts/1111",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(1111))).Times(1).Return(account, sql.ErrNoRows)
			},
			checker: func(recorder *httptest.ResponseRecorder, err error) {
				require.True(t, db.IsError(err))
				de := db.GetError(err)
				require.NotEmpty(t, de)
				require.Equal(t, http.StatusNotFound, de.Status)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)

		tc.buildStubs(store)
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, tc.url, nil)
		require.NoError(t, err)

		handler := accountgrp.New(store)
		err = handler.Get(ctx, recorder, request)

		tc.checker(recorder, err)
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, expected Account) {
	got := getResponseData[Account](t, body)
	require.Equal(t, got, expected)
}

func randomAccount() db_generated.Account {
	return db_generated.Account{
		ID:       1,
		Owner:    random.RandomOwner(),
		Currency: random.RandomCurrency(),
		Balance:  random.RandomMoney(),
	}
}
