package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/web"
)

// Get retrieves an account by its ID, and checks if the account belongs to the user from token payload
func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	account, err := h.store.GetAccount(ctx, int64(aid))

	if err != nil {
		return db.NewError(err)
	}
	if account.Owner != payload.Username {
		return auth.NewAuthError("account does not belong to user %s", payload.Username)
	}

	return web.RespondJson(ctx, w, account, http.StatusOK)
}
