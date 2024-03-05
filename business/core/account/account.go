package account

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

func (a *Core) Create(ctx context.Context, na NewAccount) (Account, error) {
	account, err := a.store.CreateAccount(ctx, CreateAccountParams{
		Owner:    na.Owner,
		Currency: na.Currency,
		Name:     na.Name,
		Balance:  na.Balance,
	})

	if err != nil {
		return Account{}, db.NewError(err)
	}

	return account, nil
}

func (a *Core) QueryById(ctx context.Context, id int64) (Account, error) {
	account, err := a.store.GetAccount(ctx, id)
	if err != nil {
		return Account{}, db.NewError(err)
	}

	return account, nil
}

func (a *Core) Query(
	ctx context.Context,
	filter QueryFilter,
	limiter QueryLimiter,
) ([]Account, error) {
	owner := db.NewNullString(filter.Owner)
	params := ListAccountsParams{
		Owner:  owner,
		Limit:  limiter.PageSize,
		Offset: (limiter.PageId - 1) * limiter.PageSize,
	}
	accounts, err := a.store.ListAccounts(ctx, params)
	if err != nil {
		return nil, db.NewError(err)
	}

	return accounts, nil
}
