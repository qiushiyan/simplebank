package entry

import (
	"context"

	"github.com/qiushiyan/simplebank/business/data/limit"
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
	limiter limit.Limiter,
) ([]Entry, error) {
	params := ListEntriesParams{
		AccountID: db.NewInt8(filter.AccountId),
		StartDate: db.NewTimestamp(filter.StartDate),
		EndDate:   db.NewTimestamp(filter.EndDate),
		Limit:     limiter.Limit,
		Offset:    limiter.Offset,
	}
	entries, err := c.store.ListEntries(ctx, params)
	if err != nil {
		return nil, db.NewError(err)
	}

	return entries, nil
}
