package friend

import (
	"errors"
)

type QueryFilter struct {
	FromAccountId *int64
	ToAccountId   *int64
	Status        *Status
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

func (q *QueryFilter) WithStatus(val Status) {
	q.Status = &val
}

func (q *QueryFilter) Valid() error {
	if q.FromAccountId == nil && q.ToAccountId == nil {
		return errors.New("must have from_account_id or to_account_id")
	}

	if q.FromAccountId != nil && q.ToAccountId != nil {
		return errors.New("cannot have both from_account_id and to_account_id")
	}

	return nil
}
