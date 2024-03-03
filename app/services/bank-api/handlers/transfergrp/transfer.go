package transfergrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type TransferRequest struct {
	FromAccountID int64 `json:"from_account_id" validate:"required"`
	ToAccountId   int64 `json:"to_account_id"   validate:"required"`
	Amount        int64 `json:"amount"          validate:"required,gt=0"`
}

func (h *Handler) Transfer(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req TransferRequest
	err := web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	if err := validate.Check(req); err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	fromAccount, err := h.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		return db.NewError(err)
	}

	if fromAccount.Owner != payload.Username {
		return auth.NewAuthError("account does not belong to user %s", payload.Username)
	}

	toAccount, err := h.store.GetAccount(ctx, req.ToAccountId)
	if err != nil {
		return db.NewError(err)
	}

	if fromAccount.Currency != toAccount.Currency {
		return fmt.Errorf("currency mismatch: %s, %s", fromAccount.Currency, toAccount.Currency)
	}

	if fromAccount.Balance < req.Amount {
		return fmt.Errorf("insufficient balance: %d", fromAccount.Balance)
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := h.store.TransferTx(ctx, args)
	if err != nil {
		return db.NewError(err)
	}

	return web.RespondJson(ctx, w, result, http.StatusCreated)
}
