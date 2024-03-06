package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type UpdateNameRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handler) UpdateName(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	var req UpdateNameRequest
	err = web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	account, err := h.core.QueryById(ctx, int64(aid))
	if err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	if payload.Username != account.Owner {
		return auth.NewForbiddenError(payload.Username)
	}

	account, err = h.core.UpdateName(ctx, account.ID, req.Name)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, account, http.StatusOK)
}