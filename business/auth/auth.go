// Package auth provides authentication and authorization support.
package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/qiushiyan/simplebank/business/auth/token"
)

func Authenticate(ctx context.Context, bearerToken string) (token.Payload, error) {
	if bearerToken == "" {
		return token.Payload{}, ErrUnauthenticated
	}

	parts := strings.Fields(bearerToken)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return token.Payload{}, errors.New("expected authorization header format: Bearer <token>")
	}

	payload, err := token.Decrypt(parts[1])
	if err != nil {
		return token.Payload{}, err
	}

	return payload, nil
}

type ctxKey int

var payloadKey ctxKey = 1

func SetPayload(ctx context.Context, p token.Payload) context.Context {
	return context.WithValue(ctx, payloadKey, p)
}

func GetPayload(ctx context.Context) token.Payload {
	val := ctx.Value(payloadKey)
	if val == nil {
		return token.Payload{}
	}
	return val.(token.Payload)
}
