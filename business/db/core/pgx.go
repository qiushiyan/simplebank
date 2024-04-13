package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Functions below are used to convert go types to a nullable sql field
// NewText converts a string pointer to a pgtype.Text
func NewText(s *string) pgtype.Text {
	var text pgtype.Text
	if s != nil {
		text.Valid = true
		text.String = *s

	}
	return text
}

// NewInt8 converts an int64 pointer to a pgtype.Int8
func NewInt8(i *int64) pgtype.Int8 {
	var i8 pgtype.Int8
	if i != nil {
		i8.Valid = true
		i8.Int64 = *i
	}
	return i8
}

// NewTimestamp converts a time.Time pointer to a pgtype.Timestamp
func NewTimestamp(t *time.Time) pgtype.Timestamp {
	var ts pgtype.Timestamp
	if t != nil {
		ts.Valid = true
		ts.Time = *t
	}
	return ts
}

// NewBool converts a bool pointer to a pgtype.Bool
func NewBool(b *bool) pgtype.Bool {
	var pb pgtype.Bool
	if b != nil {
		pb.Valid = true
		pb.Bool = *b
	}

	return pb
}
