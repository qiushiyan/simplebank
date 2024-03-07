package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/entrygrp"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetEntries(t *testing.T) {
	url := "/entries"

	cases := []struct {
		name       string
		req        entrygrp.ListEntriesRequest
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: userToken,
			req: entrygrp.ListEntriesRequest{
				FromAccountId: userAccountId,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), userAccountId).
					Times(1).
					Return(db_generated.Account(userAccount), nil)

				store.EXPECT().
					ListEntries(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db_generated.Entry{}, nil)

			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			// list entries without token
			name:  "unauthenticated",
			token: "",
			req: entrygrp.ListEntriesRequest{
				FromAccountId: userAccountId,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), userAccountId).
					Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			// list entries of a different user
			name:  "forbidden",
			token: userToken,
			req: entrygrp.ListEntriesRequest{
				FromAccountId: adminAccountId,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), adminAccountId).
					Times(1).
					Return(db_generated.Account(adminAccount), nil)

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).Times(0)
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

			body, err := json.Marshal(tc.req)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(body))
			require.NoError(t, err)
			request.Header.Add("authorization", tc.token)

			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})
	}
}
