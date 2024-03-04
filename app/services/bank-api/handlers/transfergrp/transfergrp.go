package transfergrp

import (
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/core/transfer"
	db "github.com/qiushiyan/simplebank/business/db/core"
)

type Handler struct {
	transferCore transfer.Core
	accountCore  account.Core
}

func New(store db.Store) Handler {
	return Handler{
		transferCore: transfer.NewCore(store),
		accountCore:  account.NewCore(store),
	}
}
