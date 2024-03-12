package entrygrp

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/entry"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListEntriesQuery struct {
	FromAccountId int64
	EndDate       *time.Time
	StartDate     *time.Time
}

// List entries for an account
// pass account id in post request body
// accepts 5 query parameters, from_account_id, start_date, end_date and page_id, page_size
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var fromAccountId int64
	if aid := r.URL.Query().Get("from_account_id"); aid != "" {
		aid, err := strconv.Atoi(aid)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
		fromAccountId = int64(aid)
	} else {
		return web.NewError(errors.New("from_account_id is a required query parameter"), http.StatusBadRequest)
	}

	username := auth.GetUsername(ctx)

	account, err := h.accountCore.QueryById(ctx, fromAccountId)
	if err != nil {
		return err
	}

	if account.Owner != username {
		return auth.NewForbiddenError(username)
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
	filter.WithAccountId(fromAccountId)
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
