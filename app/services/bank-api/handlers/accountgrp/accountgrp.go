package accountgrp

import (
	"github.com/qiushiyan/simplebank/business/auth/token"
	"github.com/qiushiyan/simplebank/business/core/account"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	core account.Core
}

func New(store db.Store) *Handler {
	return &Handler{
		core: account.NewCore(store),
	}
}

func (h *Handler) Register(app *web.App) {
	app.GET(
		"/accounts/all",
		h.ListAll,
		middleware.Authenticate(),
		middleware.Authorize(token.RoleAdmin),
	)
	app.GET("/accounts", h.List, middleware.Authenticate())
	app.GET("/accounts/:id", h.Get, middleware.Authenticate())
	app.POST(
		"/accounts/:id",
		h.UpdateName,
		middleware.Authenticate(),
	)
	app.GET("/accounts/search", h.Search, middleware.Authenticate())
	app.POST("/accounts", h.Create, middleware.Authenticate())
}

var _ web.RouteGroup = (*Handler)(nil)
