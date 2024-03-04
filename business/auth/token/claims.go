package token

// Claims is the additional data stored in the paseto token, accessed via token.Get("data")
type Claims struct {
	Username string `json:"username"`
	Roles    []Role `json:"roles"`
}

// NewData creates a new token data
func NewClaims(username string, roles []Role) Claims {
	return Claims{Username: username, Roles: roles}
}
