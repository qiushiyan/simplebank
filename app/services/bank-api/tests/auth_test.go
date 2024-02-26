package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/authgrp"
	"github.com/qiushiyan/simplebank/business/auth"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSignupAPi(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	user, password := randomUser()
	args := db_generated.CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
	}
	body, err := json.Marshal(authgrp.SignupRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: password,
	})
	require.NoError(t, err)

	store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(
		args,
		password,
	)).Times(1).Return(user, nil)

	recorder := httptest.NewRecorder()
	url := "/signup"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)

	handler := authgrp.New(store)
	err = handler.Signup(ctx, recorder, request)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, recorder.Code)

	requireBodyMatchUser(t, recorder.Body, user)
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
	password = "test-test"
	email := "test@gmail.com"
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
	gotUser := getResponseData[authgrp.SignupResponse](
		t,
		body,
	)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
}
