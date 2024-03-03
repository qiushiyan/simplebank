package token

import "encoding/json"

// Claims is the additional data stored in the paseto token, accessed via token.Get("data")
type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// NewData creates a new token data
func NewClaims(username string, roles []string) Claims {
	return Claims{Username: username, Roles: roles}
}

func encodeClaims(c Claims) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func decodeClaims(s string) (Claims, error) {
	var c Claims
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return Claims{}, err
	}
	return c, nil
}
