// Package auth provides authentication and authorization support.
package auth

import (
	"errors"
)

// ErrForbidden is returned when a auth issue is identified.
var ErrForbidden = errors.New("attempted action is not allowed")
