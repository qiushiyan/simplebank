package tests

import (
	"testing"

	"github.com/qiushiyan/simplebank/business/auth/token"
	"github.com/stretchr/testify/require"
)

func TestTokenRole(t *testing.T) {
	tk, err := token.NewToken("user", []string{"ADMIN"}, 0)
	require.NoError(t, err)

	payload, err := token.Decrypt(tk.GetToken())
	require.NoError(t, err)
	require.True(t, payload.HasRole("ADMIN"))
}

func TestDecryptToken(t *testing.T) {
	username := "user"
	roles := []string{"ADMIN"}

	tk, err := token.NewToken(username, roles, 0)
	require.NoError(t, err)

	payload, err := token.Decrypt(tk.GetToken())
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
	require.Equal(t, roles, payload.Roles)
}
