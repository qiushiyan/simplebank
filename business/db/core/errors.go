package db

import (
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Error struct {
	Err    error
	Status int
}

func NewError(err error) error {
	var status int

	switch {
	case isPgError(err):
		de := err.(*pgconn.PgError)
		switch de.Code {
		case pgerrcode.UniqueViolation, pgerrcode.ForeignKeyViolation:
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
	var e *pgconn.PgError
	return errors.As(err, &e)
}

func isNoRowsError(err error) bool {
	return err == pgx.ErrNoRows
}
