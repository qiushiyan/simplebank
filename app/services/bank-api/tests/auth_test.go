package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/authgrp"
	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSignupAPi(t *testing.T) {
	user, password := randomUser()
	url := "/signup"

	cases := []struct {
		name       string
		args       authgrp.SignupRequest
		buildStubs func(*mockdb.MockStore)
		checker    func(*httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			args: authgrp.SignupRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db_generated.CreateUserParams{
					Username:       user.Username,
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(
					arg,
					password,
				)).Times(1).Return(user, nil)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "DuplicatedEmail",
			args: authgrp.SignupRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db_generated.CreateUserParams{
					Username:       user.Username,
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(
					arg,
					password,
				)).Times(1).Return(db_generated.User{}, &db.ErrUniqueViolation)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			args: authgrp.SignupRequest{
				Username: user.Username,
				Email:    "invalid-email",
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			args: authgrp.SignupRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: "123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			args: authgrp.SignupRequest{
				Username: "aa",
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
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

			body, err := json.Marshal(tc.args)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			recorder := serveRequest(t, request, tc.buildStubs)

			tc.checker(recorder)
		})
	}

}

type eqCreateUserParamsMatcher struct {
	arg      db_generated.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db_generated.CreateUserParams)
	if !ok {
		return false
	}

	ok = auth.VerifyPassword(arg.HashedPassword, e.password)
	if !ok {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// EqCreateUserParams verifies that the hashedPassword field of the CreateUserParams is the hashed version of the password field in the request body
func EqCreateUserParams(arg db_generated.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser() (user db_generated.User, password string) {
	password = random.RandomPassword()
	email := random.RandomEmail()
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal("failed to hash password")
	}

	user = db_generated.User{
		Username:       random.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          email,
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db_generated.User) {
	got := getResponseData[authgrp.SignupResponse](
		t,
		body,
	)

	require.Equal(t, user.Username, got.User.Username)
	require.Equal(t, user.Email, got.User.Email)
}
