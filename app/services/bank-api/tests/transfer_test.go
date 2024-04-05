package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/transfergrp"
	db "github.com/qiushiyan/simplebank/business/db/core"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferApi(t *testing.T) {

	cases := []struct {
		name       string
		token      string
		args       db.TransferTxParams
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: userToken,
			args: db.TransferTxParams{
				FromAccountId: userAccountId,
				ToAccountId:   adminAccountId,
				Amount:        10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), db.TransferTxParams{
						FromAccountId: userAccountId,
						ToAccountId:   adminAccountId,
						Amount:        10,
					}).
					Times(1).
					Return(db.TransferTxResult{}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), userAccountId).
					Times(1).
					Return(userAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), adminAccountId).
					Times(1).
					Return(adminAccount, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body := transfergrp.TransferRequest{
				ToAccountId:   tc.args.ToAccountId,
				FromAccountId: tc.args.FromAccountId,
				Amount:        tc.args.Amount,
			}

			b, err := json.Marshal(body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(b))
			require.NoError(t, err)

			request.Header.Set("Authorization", tc.token)

			recorder := serveRequest(t, request, tc.buildStubs)

			tc.checker(recorder)
		})
	}
}
