package tests

import (
	"testing"

	"github.com/qiushiyan/simplebank/business/auth/token"
	"github.com/stretchr/testify/require"
)

func TestTokenRole(t *testing.T) {
	tk, err := token.NewToken("user", []token.Role{token.RoleAdmin}, 0)
	require.NoError(t, err)

	payload, err := token.Decrypt(tk.GetToken())
	require.NoError(t, err)
	require.True(t, payload.HasRole(token.RoleAdmin))
	require.False(t, payload.HasRole(token.RoleUser))
}

func TestDecryptToken(t *testing.T) {
	username := "user"
	roles := []token.Role{token.RoleAdmin}

	tk, err := token.NewToken(username, roles, 0)
	require.NoError(t, err)

	payload, err := token.Decrypt(tk.GetToken())
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
	require.Equal(t, roles, payload.Roles)
}
