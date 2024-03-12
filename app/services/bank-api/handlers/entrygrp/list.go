package entrygrp

import (
	"context"
	"net/http"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/entry"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListEntriesRequest struct {
	FromAccountId int64 `json:"from_account_id" validate:"required"`
}

type ListEntriesQuery struct {
	EndDate   *time.Time
	StartDate *time.Time
}

// List entries for an account
// pass account id in post request body
// accepts 4 query parameters, start_date, end_date and page_id, page_size
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req ListEntriesRequest
	err := web.Decode(r, &req)
	if err != nil {
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
		return auth.NewForbiddenError(payload.Username)
	}

	var q ListEntriesQuery

	if r.FormValue("start_date") != "" {
		val, err := time.Parse(time.DateOnly, r.FormValue("start_date"))
		q.StartDate = &val
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if r.FormValue("end_date") != "" {
		val, err := time.Parse(time.DateOnly, r.FormValue("end_date"))
		q.EndDate = &val
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if err := validate.Check(q); err != nil {
		return err
	}

	filter := entry.NewQueryFilter()
	filter.WithAccountId(req.FromAccountId)
	filter.WithEndDate(q.EndDate)
	filter.WithStartDate(q.StartDate)

	limiter, err := limit.Parse(r, 1, 20)
	if err != nil {
		return err
	}

	entries, err := h.entryCore.Query(
		ctx,
		filter,
		limiter,
	)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, entries, http.StatusOK)
}
