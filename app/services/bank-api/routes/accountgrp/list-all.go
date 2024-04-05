package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListAllAccountsResponse struct {
	Data []db_generated.Account `json:"data"`
}

// ListAllAccounts godoc
//
//	@Summary		List all accounts
//	@Description	list all accounts, available only to admin
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page_id		query	int	false	"page id, default to 1"
//	@Param			page_size	query	int	false	"page size, default to 50"
//
//	@Security		Bearer
//
//	@Success		200	{object}	ListAllAccountsResponse
//	@Failure		401
//	@Failure		409
//	@Router			/accounts/all [get]
//
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
