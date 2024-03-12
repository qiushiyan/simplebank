package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/web"
)

// List accounts for a user
// accepts two query parameters: page_id and page_size
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	filter := account.NewQueryFilter()
	filter.WithOwner(payload.Username)

	limiter, err := limit.Parse(r, 1, 5)
	if err != nil {
		return err
	}

	accounts, err := h.core.Query(
		ctx,
		filter,
		limiter,
	)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
