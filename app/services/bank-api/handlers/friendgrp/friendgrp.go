package friendgrp

import (
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/friend"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	core friend.Core
}

func New(store db.Store) *Handler {
	return &Handler{
		core: friend.NewCore(store),
	}
}

func (h Handler) Register(a *web.App) {
	a.Handle(http.MethodGet, "/friend/create", h.Create)
	a.Handle(http.MethodGet, "/friend/list", h.List)
}
