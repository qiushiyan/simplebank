package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	. "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	handler := accountgrp.New(store)
	err = handler.Get(ctx, recorder, request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func randomAccount() Account {
	return Account{
		ID:       1,
		Owner:    random.RandomOwner(),
		Currency: random.RandomCurrency(),
		Balance:  random.RandomMoney(),
	}
}
