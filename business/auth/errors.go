package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth/token"
)

var (
	ErrUnauthenticated = NewAuthError(
		"missing token in `authorization` header",
		http.StatusUnauthorized,
	)
)

// AuthError is used to pass an error during the request through the
// application with auth specific context.
type AuthError struct {
	Status int
	msg    string
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(msg string, status int) error {
	return &AuthError{
		msg:    msg,
		Status: status,
	}
}

func NewUnauthorizedError(need token.Role, args ...any) error {
	return NewAuthError(
		fmt.Sprintf("attempted action now allowed, need %s", need.Name()),
		http.StatusUnauthorized,
	)
}

func NewForbiddenError(username string) error {
	return NewAuthError(
		fmt.Sprintf("account does not belong to user %s", username),
		http.StatusForbidden,
	)
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (ae *AuthError) Error() string {
	return ae.msg
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}

func GetAuthError(err error) *AuthError {
	var ae *AuthError
	if errors.As(err, &ae) {
		return ae
	}
	return nil
}
