package entrygrp

import (
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/core/entry"
	db "github.com/qiushiyan/simplebank/business/db/core"
)

type Handler struct {
	accountCore account.Core
	entryCore   entry.Core
}

func New(store db.Store) *Handler {
	return &Handler{
		accountCore: account.NewCore(store),
		entryCore:   entry.NewCore(store),
	}
}
