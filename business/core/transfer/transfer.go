package transfer

import (
	"context"

	db "github.com/qiushiyan/simplebank/business/db/core"
)

type Core struct {
	store db.Store
}

func NewCore(store db.Store) Core {
	return Core{store: store}
}

func (t *Core) Create(ctx context.Context, nt NewTransfer) (db.TransferTxResult, error) {
	transfer, err := t.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountId: nt.FromAccountId,
		ToAccountId:   nt.ToAccountId,
		Amount:        nt.Amount,
	})
	if err != nil {
		return db.TransferTxResult{}, db.NewError(err)
	}

	return transfer, nil
}
