package accountgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/account"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type CreateAccountRequest struct {
	// less than 30 characters
	Name string `json:"name"     validate:"required,accountname"`
	// one of: USD, EUR, CAD
	Currency string `json:"currency" validate:"required,currency"`
}

type CreateAccountResponse struct {
	Data db_generated.Account `json:"data"`
}

// CreateAccount godoc
//
//	@Summary		Create an account
//	@Description	Create an account with the given name and currency. Currency should be one of "USD", "EUR" or "CAD". Name-Currency combination should be unique.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			body	body	CreateAccountRequest	true	"request body"
//
//	@Security		Bearer
//
//	@Success		200	{object}	CreateAccountResponse
//	@Failure		400
//	@Failure		409
//	@Router			/accounts/create [post]
func (h *Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateAccountRequest
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	username := auth.GetUsername(ctx)

	ret, err := h.core.Create(ctx, account.NewAccount{
		Owner:    username,
		Name:     req.Name,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, ret, http.StatusCreated)
}
