package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type GetAccountResponse struct {
	Data db_generated.Account `json:"data"`
}

// GetAccount godoc
//
//	@Summary		Get an account
//	@Description	Get account by id, token should match account owner
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"account id"
//
//	@Security		Bearer
//
//	@Success		200	{object}	GetAccountResponse
//	@Failure		401
//	@Failure		409
//	@Router			/accounts/{id} [get]
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
