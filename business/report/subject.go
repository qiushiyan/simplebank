package report

import "errors"

const (
	SubjectActivity = "activity"
)

var ErrInvalidSubject = errors.New(
	"invalid report subject, choose from 'activity'",
)

type SubjectActivityData struct {
	Username string
}
