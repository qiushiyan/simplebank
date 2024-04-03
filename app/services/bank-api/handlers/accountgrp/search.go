package accountgrp

import (
	"context"
	"errors"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SearchAccountsResponse struct {
	Data []db_generated.Account `json:"data"`
}

// SearchAccounts godoc
//
//	@Summary		Search accounts
//	@Description	Search accounts by owner
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			owner		query	string	true	"owner name"
//	@Param			page_id		query	int		false	"page id, default to 1"
//	@Param			page_size	query	int		false	"page size, default to 10"
//
//	@Security		Bearer
//
//	@Success		200	{object}	SearchAccountsResponse
//	@Failure		400
//	@Failure		409
//	@Router			/accounts/search [get]
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
