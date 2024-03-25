package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewText(s *string) pgtype.Text {
	var text pgtype.Text
	if s != nil {
		text.Valid = true
		text.String = *s

	}
	return text
}

func NewInt8(i *int64) pgtype.Int8 {
	var i8 pgtype.Int8
	if i != nil {
		i8.Valid = true
		i8.Int64 = *i
	}
	return i8
}

func NewTimestamp(t *time.Time) pgtype.Timestamp {
	var ts pgtype.Timestamp
	if t != nil {
		ts.Valid = true
		ts.Time = *t
	}
	return ts
}

func NewBool(b *bool) pgtype.Bool {
	var pb pgtype.Bool
	if b != nil {
		pb.Valid = true
		pb.Bool = *b
	}

	return pb
}
