package db

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUniqueViolation = pgconn.PgError{
		Code: pgerrcode.UniqueViolation,
	}
)

type Error struct {
	Err    error
	Status int
}

func NewError(err error) error {
	var status int

	switch {
	case isPgConnectionError(err):
		status = http.StatusServiceUnavailable
		err = fmt.Errorf("database connection error: %w", err)
	case isNoRowsError(err):
		status = http.StatusNotFound
	case isPgError(err):
		de := GetPgError(err)
		switch de.Code {
		case pgerrcode.UniqueViolation:
			status = http.StatusConflict
			switch de.ConstraintName {
			case "users_pkey", "users_email_key":
				err = errors.New(de.Detail)
			case "owner_name_key":
				err = errors.New("can't create account with the same name")
			case "owner_currency_key":
				err = errors.New("can't create account with the same currency")
			}
		case pgerrcode.ForeignKeyViolation:
			status = http.StatusConflict
		default:
			err = errors.New(de.Detail)
			status = http.StatusInternalServerError
		}

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

func isPgConnectionError(err error) bool {
	var e *pgconn.ConnectError
	return errors.As(err, &e)
}

func isPgError(err error) bool {
	var e *pgconn.PgError
	return errors.As(err, &e)
}

func GetPgError(err error) *pgconn.PgError {
	var e *pgconn.PgError
	if !errors.As(err, &e) {
		return nil
	}
	return e
}

func isNoRowsError(err error) bool {
	return err == pgx.ErrNoRows
}
