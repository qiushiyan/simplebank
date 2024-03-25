package friend

import (
	"errors"
)

type QueryFilter struct {
	FromAccountId *int64
	ToAccountId   *int64
	Pending       *bool
	Accepted      *bool
}

func NewQueryFilter() QueryFilter {
	return QueryFilter{}
}

func (q *QueryFilter) WithFromAccountID(val int64) {
	q.FromAccountId = &val
}

func (q *QueryFilter) WithToAccountID(val int64) {
	q.ToAccountId = &val
}

func (q *QueryFilter) WithPending(val *bool) {
	q.Pending = val
}

func (q *QueryFilter) WithAccepted(val *bool) {
	q.Accepted = val
}

func (q *QueryFilter) Valid() error {
	if q.FromAccountId == nil && q.ToAccountId == nil {
		return errors.New("must have at least one of from_account_id or to_account_id")
	}

	return nil
}
