package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type Payload struct {
	paseto.JSONToken
	Username string
	Roles    []string
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	}
	return nil
}

func (payload *Payload) HasRole(role string) bool {
	for _, userRole := range payload.Roles {
		if role == userRole {
			return true
		}
	}
	return false
}

// isAdmin checks if the user is an admin
func (payload *Payload) IsAdmin() bool {
	return payload.HasRole("ADMIN")
}
