package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/web"
)

// List all accounts in the database
// this is protected by the authorize middleware and can only be called by admin
// accepts two query parameters: page_id and page_size
func (h *Handler) ListAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter := account.NewQueryFilter()
	limiter, err := limit.Parse(r, 1, 50)
	if err != nil {
		return err
	}

	accounts, err := h.core.Query(ctx, filter, limiter)

	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
