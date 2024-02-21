package db

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

type Error struct {
	Err    error
	Status int
}

func NewError(err error) error {
	var status int

	switch {
	case isPgError(err):
		de := err.(*pq.Error)
		switch de.Code.Name() {
		case "unique_violation", "foreign_key_violation":
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
		}
	case isNoRowsError(err):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	return &Error{Err: err, Status: status}
}

func (de *Error) Error() string {
	if isNoRowsError(de.Err) {
		return "no rows found"
	}
	return de.Err.Error()
}

func GetError(err error) *Error {
	var de *Error
	if !errors.As(err, &de) {
		return nil
	}
	return de
}

func IsError(err error) bool {
	var re *Error
	return errors.As(err, &re)
}

func isPgError(err error) bool {
	var e *pq.Error
	return errors.As(err, &e)
}

func isNoRowsError(err error) bool {
	return err == sql.ErrNoRows
}
