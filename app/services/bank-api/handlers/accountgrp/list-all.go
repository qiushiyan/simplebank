package accountgrp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/core/account"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListAllAccountQuery struct {
	PageId   int32 `json:"page_id"   validate:"min=1"`
	PageSize int32 `json:"page_size" validate:"min=1,max=50"`
}

// List all accounts in the database
// this is protected by the authorize middleware and can only be called by admin
func (h *Handler) ListAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var q ListAllAccountQuery
	if r.FormValue("page_id") == "" {
		q.PageId = 1
	} else {
		id, err := strconv.Atoi(r.FormValue("page_id"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.PageId = int32(id)
	}

	if r.FormValue("page_size") == "" {
		q.PageSize = 50
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

	accounts, err := h.core.Query(ctx, account.QueryFilter{}, account.QueryLimiter{
		PageId:   q.PageId,
		PageSize: q.PageSize,
	})

	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, accounts, http.StatusOK)
}
