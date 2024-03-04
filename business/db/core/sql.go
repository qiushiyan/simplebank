package db

import (
	"database/sql"
	"time"
)

func NewNullString(s *string) sql.NullString {
	var ns sql.NullString
	if s != nil {
		ns.Valid = true
		ns.String = *s

	}
	return ns
}

func NewNullInt64(i *int64) sql.NullInt64 {
	var ni sql.NullInt64
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	}
	return ni
}

func NewNullTime(t *time.Time) sql.NullTime {
	var nt sql.NullTime
	if t != nil {
		nt.Valid = true
		nt.Time = *t
	}
	return nt
}
