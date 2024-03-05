package entry

import (
	"context"

	db "github.com/qiushiyan/simplebank/business/db/core"
	. "github.com/qiushiyan/simplebank/business/db/generated"
)

type Core struct {
	store db.Store
}

func NewCore(store db.Store) Core {
	return Core{store: store}
}

func (c *Core) Query(
	ctx context.Context,
	filter QueryFilter,
	limiter QueryLimiter,
) ([]Entry, error) {
	params := ListEntriesParams{
		AccountID: db.NewNullInt64(filter.AccountId),
		StartDate: db.NewNullTime(filter.StartDate),
		EndDate:   db.NewNullTime(filter.EndDate),
		Limit:     limiter.PageSize,
		Offset:    (limiter.PageId - 1) * limiter.PageSize,
	}
	entries, err := c.store.ListEntries(ctx, params)
	if err != nil {
		return nil, db.NewError(err)
	}

	return entries, nil
}