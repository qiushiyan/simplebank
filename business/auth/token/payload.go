package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type Payload struct {
	paseto.JSONToken
	Username string
	Roles    []Role
}

func (payload Payload) IsEmpty() bool {
	return payload.Username == ""
}

// Valid checks if the token payload is valid or not
func (payload Payload) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	}
	return nil
}

func (payload Payload) HasRole(role Role) bool {
	for i := range payload.Roles {
		if role == payload.Roles[i] {
			return true
		}
	}
	return false
}

// isAdmin checks if the user is an admin
func (payload Payload) IsAdmin() bool {
	return payload.HasRole(RoleAdmin)
}
