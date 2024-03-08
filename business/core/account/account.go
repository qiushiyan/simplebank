package account

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

func (a *Core) Create(ctx context.Context, na NewAccount) (Account, error) {
	account, err := a.store.CreateAccount(ctx, CreateAccountParams(na))

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
	limiter limit.Limiter,
) ([]Account, error) {
	owner := db.NewText(filter.Owner)
	params := ListAccountsParams{
		Owner:  owner,
		Limit:  limiter.Limit,
		Offset: limiter.Offset,
	}
	accounts, err := a.store.ListAccounts(ctx, params)
	if err != nil {
		return nil, db.NewError(err)
	}

	return accounts, nil
}

func (a *Core) UpdateName(ctx context.Context, id int64, name string) (Account, error) {
	account, err := a.store.UpdateAccountName(ctx, UpdateAccountNameParams{
		ID:   id,
		Name: name,
	})

	if err != nil {
		return Account{}, db.NewError(err)
	}

	return account, nil
}
