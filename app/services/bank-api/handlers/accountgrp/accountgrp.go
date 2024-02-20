package accountgrp

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	store db.Store
}

func New(store db.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := strings.Split(r.URL.Path, "/")[2]
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ret, err := h.store.GetAccount(ctx, int64(aid))
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, ret, http.StatusOK)
}

func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ret, err := h.store.ListAccounts(ctx, db_generated.ListAccountsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, ret, http.StatusOK)
}
