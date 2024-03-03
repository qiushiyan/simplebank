package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListAccountQuery struct {
	PageID   int32 `json:"page_id"   validate:"min=1"`
	PageSize int32 `json:"page_size" validate:"min=1,max=5"`
}

// List accounts for a user
// accepts two query parameters: page_id and page_size
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var q ListAccountQuery

	if r.FormValue("page_id") == "" {
		q.PageID = 1
	} else {
		id, err := strconv.Atoi(r.FormValue("page_id"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.PageID = int32(id)
	}

	if r.FormValue("page_size") == "" {
		q.PageSize = 5
	} else {
		size, err := strconv.Atoi(r.FormValue("page_size"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.PageSize = int32(size)
	}

	if err := validate.Check(q); err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	accounts, err := h.store.ListAccounts(ctx, db_generated.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  q.PageSize,
		Offset: (q.PageID - 1) * q.PageSize,
	})

	if err != nil {
		return db.NewError(err)
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
