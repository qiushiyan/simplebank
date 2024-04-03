package db

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qiushiyan/simplebank/foundation/web"
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

	// connection error
	case isPgConnectionError(err):
		status = http.StatusServiceUnavailable
		err = fmt.Errorf("database connection error: %w", err)

	// no result set
	case IsNoRowsError(err):
		status = http.StatusNotFound

	// postgres server error
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
			if de.Detail != "" {
				err = errors.New(de.Detail)
			}

			if de.Code == "42P01" {
				return web.NewShutdownError(
					fmt.Sprintf("%s, did you forget the run the migrations?", err.Error()),
				)
			}
			status = http.StatusInternalServerError
		}

	default:
		status = http.StatusInternalServerError
	}

	return &Error{Err: err, Status: status}
}

func (de *Error) Error() string {
	if IsNoRowsError(de.Err) {
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

func IsNoRowsError(err error) bool {
	return err == pgx.ErrNoRows
}
