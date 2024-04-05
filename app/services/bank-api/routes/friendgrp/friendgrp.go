// Package friendgrp provides friendship-related handlers
package friendgrp

import (
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/core/friend"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	friend  friend.Core
	account account.Core
}

func New(store db.Store) *Handler {
	return &Handler{
		friend:  friend.NewCore(store),
		account: account.NewCore(store),
	}
}

func (h Handler) Register(a *web.App) {
	a.Handle(http.MethodGet, "/friend/create", h.Create, middleware.Authenticate())
	a.Handle(http.MethodGet, "/friend/list", h.List, middleware.Authenticate())
}
