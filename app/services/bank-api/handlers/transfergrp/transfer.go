package transfergrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/transfer"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type TransferRequest struct {
	FromAccountId int64 `json:"from_account_id" validate:"required"`
	ToAccountId   int64 `json:"to_account_id"   validate:"required"`
	Amount        int64 `json:"amount"          validate:"required,gt=0"`
}

func (h *Handler) Transfer(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req TransferRequest
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	fromAccount, err := h.accountCore.QueryById(ctx, req.FromAccountId)
	if err != nil {
		return err
	}

	if fromAccount.Owner != payload.Username {
		return auth.NewForbiddenError(payload.Username)
	}

	toAccount, err := h.accountCore.QueryById(ctx, req.ToAccountId)
	if err != nil {
		return err
	}

	if fromAccount.Currency != toAccount.Currency {
		return fmt.Errorf("currency mismatch: %s, %s", fromAccount.Currency, toAccount.Currency)
	}

	if fromAccount.Balance < req.Amount {
		return fmt.Errorf("insufficient balance: %d", fromAccount.Balance)
	}

	result, err := h.transferCore.Create(ctx, transfer.NewTransfer(req))
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, result, http.StatusCreated)
}
