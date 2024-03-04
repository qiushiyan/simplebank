package entry

import "time"

type QueryFilter struct {
	AccountId *int64
	StartDate *time.Time
	EndDate   *time.Time
}

type QueryLimiter struct {
	PageId   int32
	PageSize int32
}
