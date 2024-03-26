package accountgrp

import (
	"context"
	"errors"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/web"
)

func (h *Handler) Search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Query().Get("owner") == "" {
		return web.NewError(
			errors.New("`owner` is a required query parameter"),
			http.StatusBadRequest,
		)
	}

	owner := r.URL.Query().Get("owner")
	filter := account.NewQueryFilter()
	filter.WithOwner(owner)

	limiter, err := limit.Parse(r, 1, 10)
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	accounts, err := h.core.Query(ctx, filter, limiter)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
