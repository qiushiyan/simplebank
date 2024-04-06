package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/friendgrp"
	"github.com/qiushiyan/simplebank/business/core/friend"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateFriend(t *testing.T) {
	url := "/friend/create"
	fromAccountId := int64(1)
	toAccountId := int64(2)
	arg := friendgrp.CreateFriendRequest{
		FromAccountId: fromAccountId,
		ToAccountId:   toAccountId,
	}

	cases := []struct {
		name       string
		token      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: userToken,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountId).
					Times(1).
					Return(db_generated.Account{Owner: "user"}, nil)

				store.EXPECT().
					CreateFriend(gomock.Any(), db_generated.CreateFriendParams{
						FromAccountID: fromAccountId,
						ToAccountID:   toAccountId,
					}).
					Times(1).
					Return(db_generated.Friendship{
						FromAccountID: fromAccountId,
						ToAccountID:   toAccountId,
						Status:        "pending",
					}, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				friendship := getResponseData[db_generated.Friendship](t, recorder.Body)
				require.Equal(t, fromAccountId, friendship.FromAccountID)
				require.Equal(t, toAccountId, friendship.ToAccountID)
				require.Equal(t, "pending", friendship.Status)
			},
		},
		{
			name:  "unauthenticated",
			token: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), fromAccountId).Times(0)
				store.EXPECT().CreateFriend(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "wrong-owner",
			token: userToken,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountId).
					Times(1).
					Return(db_generated.Account{Owner: "admin"}, nil)
				store.EXPECT().CreateFriend(gomock.Any(), gomock.Any()).Times(0)
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
			body, err := json.Marshal(arg)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))

			request.Header.Set("authorization", tc.token)
			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})
	}
}

func TestListFriends(t *testing.T) {
	url := "/friend/list"

	cases := []struct {
		name       string
		token      string
		query      string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:  "ok-from",
			token: userToken,
			query: "from_account_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				fromAccountId := int64(1)

				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountId).
					Times(1).
					Return(db_generated.Account{Owner: "user"}, nil)

				store.EXPECT().ListFriends(gomock.Any(), db_generated.ListFriendsParams{
					FromAccountID: db.NewInt8(&fromAccountId),
					Offset:        0,
					Limit:         50,
				}).Times(1).Return([]db_generated.Friendship{
					{
						FromAccountID: 1,
						ToAccountID:   2,
						Status:        "accepted",
					},
				}, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				friendships := getResponseData[[]db_generated.Friendship](t, recorder.Body)
				require.Len(t, friendships, 1)
				require.Equal(t, int64(1), friendships[0].FromAccountID)
				require.Equal(t, int64(2), friendships[0].ToAccountID)
				require.Equal(t, "accepted", friendships[0].Status)
			},
		},
		{
			name:  "ok-to",
			token: userToken,
			query: "to_account_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				toAccountId := int64(1)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccountId).
					Times(1).
					Return(db_generated.Account{Owner: "user"}, nil)

				store.EXPECT().ListFriends(gomock.Any(), db_generated.ListFriendsParams{
					ToAccountID: db.NewInt8(&toAccountId),
					Offset:      0,
					Limit:       50,
				}).Times(1).Return([]db_generated.Friendship{
					{
						FromAccountID: 1,
						ToAccountID:   2,
						Status:        "accepted",
					},
				}, nil)

			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				friendships := getResponseData[[]db_generated.Friendship](t, recorder.Body)
				require.Len(t, friendships, 1)
				require.Equal(t, int64(1), friendships[0].FromAccountID)
				require.Equal(t, int64(2), friendships[0].ToAccountID)
				require.Equal(t, "accepted", friendships[0].Status)
			},
		},
		{
			name:  "ok-status",
			token: userToken,
			query: "from_account_id=1&status=accepted",
			buildStubs: func(store *mockdb.MockStore) {
				fromAccountId := int64(1)
				status := "accepted"

				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountId).
					Times(1).
					Return(db_generated.Account{Owner: "user"}, nil)

				store.EXPECT().ListFriends(gomock.Any(), db_generated.ListFriendsParams{
					FromAccountID: db.NewInt8(&fromAccountId),
					Status:        db.NewText(&status),
					Offset:        0,
					Limit:         50,
				}).Times(1).Return([]db_generated.Friendship{
					{
						FromAccountID: 1,
						ToAccountID:   2,
						Status:        "accepted",
					},
				}, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				friendships := getResponseData[[]db_generated.Friendship](t, recorder.Body)
				require.Len(t, friendships, 1)
				require.Equal(t, int64(1), friendships[0].FromAccountID)
				require.Equal(t, int64(2), friendships[0].ToAccountID)
				require.Equal(t, "accepted", friendships[0].Status)
			},
		},
		{
			name:  "unauthenticated",
			token: "",
			query: "from_account_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().ListFriends(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "wrong-owner",
			token: userToken,
			query: "from_account_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), int64(1)).Return(db_generated.Account{
					Owner: "admin",
				}, nil).Times(1)

				store.EXPECT().ListFriends(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:  "invalid-status",
			token: userToken,
			query: "status=invalid-status",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().ListFriends(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "invalid-both-from-to",
			token: userToken,
			query: "from_account_id=1&to_account_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().ListFriends(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			request.Header.Set("authorization", tc.token)
			request.URL.RawQuery = tc.query

			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})
	}
}

func TestUpdateFriend(t *testing.T) {
	url := "/friend/1"

	cases := []struct {
		name       string
		token      string
		status     string
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name:   "ok",
			token:  userToken,
			status: friend.StatusAccepted.Name(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFriend(gomock.Any(), int64(1)).
					Times(1).
					Return(db_generated.Friendship{ToAccountID: 1}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), int64(1)).
					Times(1).
					Return(db_generated.Account{Owner: "user"}, nil)

				store.EXPECT().UpdateFriend(gomock.Any(), db_generated.UpdateFriendParams{
					ID:     1,
					Status: "accepted",
				}).Times(1).Return(db_generated.Friendship{
					Status: "accepted",
				}, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				friendship := getResponseData[db_generated.Friendship](t, recorder.Body)
				require.Equal(t, "accepted", friendship.Status)
			},
		},
		{
			name:  "unauthenticated",
			token: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFriend(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().UpdateFriend(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		}, {
			name:   "wrong-owner",
			token:  userToken,
			status: friend.StatusAccepted.Name(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFriend(gomock.Any(), int64(1)).
					Times(1).
					Return(db_generated.Friendship{
						ToAccountID: 1,
					}, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), int64(1)).
					Times(1).
					Return(db_generated.Account{
						Owner: "admin",
					}, nil)

				store.EXPECT().UpdateFriend(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		}, {
			name:   "invalid-status",
			token:  userToken,
			status: "invalid-status",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFriend(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().UpdateFriend(gomock.Any(), gomock.Any()).Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(friendgrp.UpdateRequest{Status: tc.status})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
			require.NoError(t, err)
			request.Header.Set("authorization", tc.token)

			recorder := serveRequest(t, request, tc.buildStubs)
			tc.checker(recorder)
		})
	}

}
