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

func Decrypt(t string) (Payload, error) {
	var payload Payload
	err := v2.Decrypt(t, key, &payload, nil)
	if err != nil {
		return Payload{}, ErrInvalidToken
	}

	s := payload.Get("data")
	c, err := decodeClaims(s)

	if err != nil {
		return Payload{}, ErrInvalidToken
	}

	payload.Username = c.Username
	payload.Roles = c.Roles

	return payload, nil
}

// NewToken creates a new paseto token with a username, a set of roles and a duration
func NewToken(username string, roles []string, duration time.Duration) (Token, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Token{}, err
	}

	var d time.Duration
	if duration == 0 {
		d = 7 * 24 * time.Hour
	}

	payload := Payload{
		JSONToken: paseto.JSONToken{
			Jti:        id.String(),
			Subject:    username,
			Issuer:     "simplebank",
			Audience:   "simplebank-user",
			IssuedAt:   time.Now(),
			Expiration: time.Now().Add(d),
		},
		Username: username,
		Roles:    roles,
	}

	s, err := encodeClaims(NewClaims(username, roles))
	if err != nil {
		return Token{}, err
	}
	payload.Set("data", s)

	token, err := v2.Encrypt(key, payload, nil)
	if err != nil {
		return Token{}, err
	}

	return Token{t: token}, nil
}

func (t Token) GetToken() string {
	return t.t
}
