package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type UpdateRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateResponse struct {
	Data db_generated.Account `json:"data"`
}

// Update godoc
//
//	@Summary		Update an account
//	@Description	Update account by id, token should match account owner, currently only name can be updated
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int				true	"account id"
//	@Param			body	body	UpdateRequest	true	"request body"
//
//	@Security		Bearer
//
//	@Success		200	{object}	UpdateResponse
//	@Failure		400
//	@Failure		403
//
//	@Router			/accounts/{id} [patch]
func (h *Handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := cast.ToInt64E(web.Param(r, "id"))
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	var req UpdateRequest
	err = web.Decode(r, &req)
	if err != nil {
		return err
	}

	account, err := h.core.QueryById(ctx, id)
	if err != nil {
		return err
	}

	username := auth.GetUsername(ctx)

	if account.Owner != username {
		return auth.NewForbiddenError(username)
	}

	account, err = h.core.UpdateName(ctx, account.ID, req.Name)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, account, http.StatusOK)
}
