package accountgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/data/limit"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListAccountsResponse struct {
	Data []db_generated.Account `json:"data"`
}

// ListAccount godoc
//
//	@Summary		List accounts for user
//	@Description	List accounts for the authenticated user in the token
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page_id		query	int	false	"page id, default to 1"
//	@Param			page_size	query	int	false	"page size, default to 5"
//	@Security		Bearer
//	@Success		200	{object}	[]db_generated.Account
//	@Failure		401
//	@Failure		409
//	@Router			/accounts [get]
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	username := auth.GetUsername(ctx)

	filter := account.NewQueryFilter()
	filter.WithOwner(username)

	limiter, err := limit.Parse(r, 1, 5)
	if err != nil {
		return err
	}

	accounts, err := h.core.Query(
		ctx,
		filter,
		limiter,
	)
	fmt.Println("accounts", accounts)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
