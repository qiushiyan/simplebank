package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/foundation/web"
)

// Get retrieves an account by its ID, and checks if the account belongs to the user from token payload
func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	username := auth.GetUsername(ctx)

	account, err := h.core.QueryById(ctx, int64(aid))
	if err != nil {
		return err
	}
	if account.Owner != username {
		return auth.NewForbiddenError(username)
	}

	return web.RespondJson(ctx, w, account, http.StatusOK)
}
