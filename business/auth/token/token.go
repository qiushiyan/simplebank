package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var key = []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
var v2 = paseto.NewV2()

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token invalid")
	ErrExpiredToken = errors.New("token expired")
)

type Token struct {
	t string
}

// NewToken creates a new token with a username, a set of roles and a duration
func NewToken(username string, roles []string, duration time.Duration) (*Token, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	var d time.Duration
	if duration == 0 {
		d = 7 * 24 * time.Hour
	}

	payload := &Payload{
		JSONToken: paseto.JSONToken{
			Jti:        id.String(),
			Subject:    username,
			Issuer:     "simplebank",
			Audience:   "simplebank-user",
			IssuedAt:   time.Now(),
			Expiration: time.Now().Add(d),
		},
	}

	s, err := encodeClaims(NewClaims(roles))
	if err != nil {
		return nil, err
	}
	payload.Set("data", s)

	token, err := v2.Encrypt(key, payload, nil)
	if err != nil {
		return nil, err
	}

	return &Token{t: token}, nil
}

func (t *Token) decrypt() (*Payload, error) {
	var payload Payload
	err := v2.Decrypt(t.t, key, &payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &payload, nil
}

func (t *Token) IsAdmin() bool {
	payload, err := t.decrypt()
	if err != nil {
		return false
	}
	if err := payload.Valid(); err != nil {
		return false
	}

	s := payload.Get("data")
	c, err := decodeClaims(s)
	if err != nil {
		return false
	}
	return c.isAdmin()
}

func (t *Token) GetToken() string {
	return t.t
}
