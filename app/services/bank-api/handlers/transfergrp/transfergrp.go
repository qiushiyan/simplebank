package transfergrp

import (
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/core/transfer"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
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

func (h Handler) Register(app *web.App) {
	app.POST("/transfer", h.Transfer, middleware.Authenticate())
}

var _ web.RouteGroup = (*Handler)(nil)
