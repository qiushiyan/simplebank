package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/authgrp"
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
	response := db.CreateUserTxResult{
		User:              user,
		AfterCreateResult: "",
	}

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
				Email:    user.Email.String,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				params := makeCreateUserTxParams(
					user.Username,
					user.Email,
					user.HashedPassword,
				)
				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserTxParams(
					params,
					password,
				)).Times(1).Return(response, nil)

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
				Email:    user.Email.String,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				params := makeCreateUserTxParams(
					user.Username,
					user.Email,
					user.HashedPassword,
				)
				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserTxParams(
					params,
					password,
				)).Times(1).Return(db.CreateUserTxResult{}, &db.ErrUniqueViolation)
			},
			checker: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			args: authgrp.SignupRequest{
				Username: user.Username,
				Email:    user.Email.String,
				Password: "123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
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
				Email:    user.Email.String,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
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

// EqCreateUserTxParams verifies that the hashedPassword field of the CreateUserParams is the hashed version of the password field in the request body
func EqCreateUserTxParams(params db.CreateUserTxParams, password string) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{params, password}
}

type eqCreateUserTxParamsMatcher struct {
	params   db.CreateUserTxParams
	password string
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	input, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	ok = auth.VerifyPassword(input.HashedPassword, e.password)
	if !ok {
		return false
	}

	e.params.HashedPassword = input.HashedPassword
	return reflect.DeepEqual(e.params.CreateUserParams, input.CreateUserParams)
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.params, e.password)
}

func randomUser() (db_generated.User, string) {
	password := random.RandomPassword()
	email := random.RandomEmail()
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		panic("failed to hash password")
	}

	user := db_generated.User{
		Username:       random.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          db.NewText(&email),
	}
	return user, password
}

func makeCreateUserTxParams(
	username string,
	email pgtype.Text,
	hashedPassword string,
) db.CreateUserTxParams {
	return db.CreateUserTxParams{
		CreateUserParams: db_generated.CreateUserParams{
			Username:       username,
			Email:          email,
			HashedPassword: hashedPassword,
		},
		AfterCreate: func(user db_generated.User) (any, error) {
			return "", nil
		},
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db_generated.User) {
	got := getResponseData[authgrp.SignupResponse](
		t,
		body,
	)

	require.Equal(t, user.Username, got.User.Username)
	require.Equal(t, user.Email.String, got.User.Email)
}
