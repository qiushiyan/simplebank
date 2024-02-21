package accountgrp

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
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

type ListAccountRequest struct {
	PageID   int32 `json:"page_id" validate:"min=1"`
	PageSize int32 `json:"page_size" validate:"min=1,max=50"`
}

func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("page_id")
	size := r.FormValue("page_size")

	if id == "" {
		id = "1"
	}

	if size == "" {
		size = "10"
	}

	pageId, err := strconv.Atoi(id)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	pageSize, err := strconv.Atoi(size)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	req := ListAccountRequest{
		PageID:   int32(pageId),
		PageSize: int32(pageSize),
	}

	if err := validate.Check(req); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	ret, err := h.store.ListAccounts(ctx, db_generated.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}

	return web.RespondJson(ctx, w, ret, http.StatusOK)
}
