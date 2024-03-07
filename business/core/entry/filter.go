package entry

import "time"

type QueryFilter struct {
	AccountId *int64
	StartDate *time.Time
	EndDate   *time.Time
}

func NewQueryFilter() QueryFilter {
	return QueryFilter{}
}

func (f *QueryFilter) WithAccountId(accountId int64) {
	f.AccountId = &accountId
}

func (f *QueryFilter) WithStartDate(startDate *time.Time) {
	f.StartDate = startDate
}

func (f *QueryFilter) WithEndDate(endDate *time.Time) {
	f.EndDate = endDate
}
