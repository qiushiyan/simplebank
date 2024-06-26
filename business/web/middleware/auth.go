package middleware

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/auth/token"
	"github.com/qiushiyan/simplebank/foundation/web"
)

// Authenticate validates a token from the `Authorization` header
// if the token is valid, the payload is added to the context
func Authenticate() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			payload, err := auth.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return err
			}

			ctx = auth.SetPayload(ctx, payload)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list.
func Authorize(role token.Role) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			payload := auth.GetPayload(ctx)
			if payload.IsEmpty() {
				return auth.ErrUnauthenticated
			}

			// check role match, skip if role is empty
			if !payload.HasRole(role) {
				return auth.NewUnauthorizedError(role, payload.Roles)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
