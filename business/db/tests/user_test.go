package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	. "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/random"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	params := CreateUserParams{
		Username:       random.RandomOwner(),
		HashedPassword: "secret",
	}

	user, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser()

	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.True(t, auth.VerifyPassword(user2.HashedPassword, "secret"))

	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func createRandomUser() User {
	hashedPassword, err := auth.HashPassword("secret")
	if err != nil {
		log.Fatal(err)
	}

	params := CreateUserParams{
		Username:       random.RandomOwner(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalf("cannot create user: %v", err)
	}
	return user
}
