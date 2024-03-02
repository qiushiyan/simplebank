package auth

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthenticated = NewAuthError("missing token in `authorization` header")
)

// AuthError is used to pass an error during the request through the
// application with auth specific context.
type AuthError struct {
	msg string
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(format string, args ...any) error {
	return &AuthError{
		msg: fmt.Sprintf(format, args...),
	}
}

func NewUnauthorizedError(need string, args ...any) error {
	return NewAuthError("attempted action now allowed, need %s", need)
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
