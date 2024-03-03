package entrygrp

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type ListEntriesRequest struct {
	FromAccountID int64 `json:"from_account_id" validate:"required"`
}

type ListEntriesQuery struct {
	PageId    int32 `json:"page_id"   validate:"min=1"`
	PageSize  int32 `json:"page_size" validate:"min=1,max=20"`
	EndDate   sql.NullTime
	StartDate sql.NullTime
}

func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req ListEntriesRequest
	err := web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	if err := validate.Check(req); err != nil {
		return err
	}

	var q ListEntriesQuery

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
		q.PageSize = 20
	} else {
		size, err := strconv.Atoi(r.FormValue("page_size"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.PageSize = int32(size)
	}

	if r.FormValue("start_date") == "" {
		q.StartDate = sql.NullTime{
			Valid: false,
		}
	} else {
		startDate, err := time.Parse(time.DateOnly, r.FormValue("start_date"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.StartDate = sql.NullTime{
			Valid: true,
			Time:  startDate,
		}
	}

	if r.FormValue("end_date") == "" {
		q.EndDate = sql.NullTime{
			Valid: false,
		}
	} else {
		endDate, err := time.Parse(time.DateOnly, r.FormValue("end_date"))
		if err != nil {
			return response.NewError(err, http.StatusBadRequest)
		}
		q.EndDate = sql.NullTime{
			Valid: true,
			Time:  endDate,
		}
	}

	if err := validate.Check(q); err != nil {
		return err
	}

	payload := auth.GetPayload(ctx)
	if payload.IsEmpty() {
		return auth.ErrUnauthenticated
	}

	account, err := h.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		return db.NewError(err)
	}

	if payload.Username != account.Owner {
		return auth.NewAuthError("account does not belong to user %s", payload.Username)
	}

	entries, err := h.store.ListEntries(ctx, db_generated.ListEntriesParams{
		AccountID: req.FromAccountID,
		Limit:     q.PageSize,
		Offset:    (q.PageId - 1) * int32(q.PageSize),
		StartDate: q.StartDate,
		EndDate:   q.EndDate,
	})

	if err != nil {
		return db.NewError(err)
	}

	return web.RespondJson(ctx, w, entries, http.StatusOK)
}
