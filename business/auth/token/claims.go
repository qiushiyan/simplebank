package token

import "encoding/json"

// Claims is the additional data stored in the paseto token, access via token.Get("data")
type Claims struct {
	Roles []string `json:"roles"`
}

// NewData creates a new token data
func NewClaims(roles []string) *Claims {
	return &Claims{Roles: roles}
}

// isAdmin checks if the user is an admin
func (c *Claims) isAdmin() bool {
	for _, role := range c.Roles {
		if role == "ADMIN" {
			return true
		}
	}
	return false
}

func encodeClaims(c *Claims) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func decodeClaims(s string) (*Claims, error) {
	var c Claims
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
