package entrygrp

import (
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/core/entry"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
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

func (h Handler) Register(app *web.App) {
	app.GET("/entries", h.List, middleware.Authenticate())
}

var _ web.RouteGroup = (*Handler)(nil)
