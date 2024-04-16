package entrygrp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/entry"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type ListEntriesQuery struct {
	// The sender account ID
	FromAccountId int64
	// Maximum date for the entry
	EndDate *time.Time
	// Minimum date for the entry
	StartDate *time.Time
}

// ListEntries godoc
//
//	@Summary		List entries for an account
//	@Description	List entries for the account from the token
//	@Tags			Entries
//	@Accept			json
//	@Produce		json
//	@Param			from_account_id	query	int		true	"Account ID"
//	@Param			start_date		query	string	false	"Start Date"
//	@Param			end_date		query	string	false	"End Date"
//	@Param			page_id			query	int		false	"Page ID"
//	@Param			page_size		query	int		false	"Page Size"
//
// @Security		Bearer
//
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/entries [get]
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	values := r.URL.Query()

	var fromAccountId int64
	var err error
	if aid := values.Get("from_account_id"); aid != "" {
		fromAccountId, err = cast.ToInt64E(aid)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
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

	if date := values.Get("start_date"); date != "" {
		val, err := time.Parse(time.DateOnly, date)
		q.StartDate = &val
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if date := r.FormValue("end_date"); date != "" {
		val, err := time.Parse(time.DateOnly, date)
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
