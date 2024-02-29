package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" validate:"required,currency"`
}

func (h *Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateAccountRequest
	err := web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	if err := validate.Check(req); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	payload := auth.GetPayload(ctx)
	if payload == nil {
		return auth.NewAuthError("missing authentication payload")
	}

	ret, err := h.store.CreateAccount(ctx, db_generated.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
	})
	if err != nil {
		return db.NewError(err)
	}

	return web.RespondJson(ctx, w, ret, http.StatusCreated)
}
