package accountgrp

import (
	"github.com/qiushiyan/simplebank/business/core/account"
	db "github.com/qiushiyan/simplebank/business/db/core"
)

type Handler struct {
	core account.Core
}

func New(store db.Store) *Handler {
	return &Handler{
		core: account.NewCore(store),
	}
}
