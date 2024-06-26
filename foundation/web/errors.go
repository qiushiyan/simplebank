package web

import "errors"

// ErrorDocument is the form used for API responses from failures in the API.
type ErrorDocument struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err    error
	Status int
}

// NewError wraps a provided error with an HTTP status code.
// This function should be used as fallback when handlers encounter expected errors
// i.e., if it was not a validation error, auth error, or db error, etc.
func NewError(err error, status int) error {
	return &Error{err, status}
}

// NewErrorS wraps the message string with errors.New and creates a new Error
func NewErrorS(msg string, status int) error {
	return &Error{Err: errors.New(msg), Status: status}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (re *Error) Error() string {
	return re.Err.Error()
}

// IsError checks if an error of type Error exists.
func IsError(err error) bool {
	var re *Error
	return errors.As(err, &re)
}

// GetError returns a copy of the Error pointer.
func GetError(err error) *Error {
	var re *Error
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
