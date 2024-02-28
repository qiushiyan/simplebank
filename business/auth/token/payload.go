package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type Payload struct {
	paseto.JSONToken
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	}
	return nil
}
