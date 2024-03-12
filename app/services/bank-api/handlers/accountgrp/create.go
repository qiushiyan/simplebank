package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type CreateAccountRequest struct {
	Name     string `json:"name"     validate:"required,accountname"`
	Currency string `json:"currency" validate:"required,currency"`
}

func (h *Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateAccountRequest
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	ret, err := h.core.Create(ctx, account.NewAccount{
		Owner:    payload.Username,
		Name:     req.Name,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, ret, http.StatusCreated)
}
