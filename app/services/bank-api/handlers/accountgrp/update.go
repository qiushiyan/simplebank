package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type UpdateNameRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handler) UpdateName(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	var req UpdateNameRequest
	err = web.Decode(r, &req)
	if err != nil {
		return err
	}

	account, err := h.core.QueryById(ctx, int64(aid))
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
