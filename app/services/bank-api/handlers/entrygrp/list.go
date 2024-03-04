package entrygrp

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/entry"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListEntriesRequest struct {
	FromAccountId int64 `json:"from_account_id" validate:"required"`
}

type ListEntriesQuery struct {
	PageId    int `json:"page_id"   validate:"min=1"`
	PageSize  int `json:"page_size" validate:"min=1,max=20"`
	EndDate   *time.Time
	StartDate *time.Time
}

// List entries for an account
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req ListEntriesRequest
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

	account, err := h.accountCore.QueryById(ctx, req.FromAccountId)
	if err != nil {
		return err
	}

	if account.Owner != payload.Username {
		return auth.NewAuthError("account does not belong to user %s", payload.Username)
	}

	var q ListEntriesQuery

	if r.FormValue("page_id") == "" {
		q.PageId = 1
	} else {
		q.PageId, err = strconv.Atoi(r.FormValue("page_id"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
	}

	if r.FormValue("page_size") == "" {
		q.PageSize = 20
	} else {
		q.PageSize, err = strconv.Atoi(r.FormValue("page_size"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
	}

	if r.FormValue("start_date") != "" {
		val, err := time.Parse(time.DateOnly, r.FormValue("start_date"))
		q.StartDate = &val
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
	}

	if r.FormValue("end_date") != "" {
		val, err := time.Parse(time.DateOnly, r.FormValue("end_date"))
		q.EndDate = &val
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
	}

	if err := validate.Check(q); err != nil {
		return err
	}

	entries, err := h.entryCore.Query(ctx, entry.QueryFilter{
		AccountId: &req.FromAccountId,
		StartDate: q.StartDate,
		EndDate:   q.EndDate,
	}, entry.QueryLimiter{
		PageId:   int32(q.PageId),
		PageSize: int32(q.PageSize),
	})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, entries, http.StatusOK)
}
