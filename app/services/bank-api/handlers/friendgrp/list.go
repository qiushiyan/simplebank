package friendgrp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/business/core/friend"
	"github.com/qiushiyan/simplebank/business/data/limit"
	"github.com/qiushiyan/simplebank/foundation/web"
)

func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var fromAccountId int
	var toAccountId int
	var pending *bool
	var accepted *bool
	var err error
	if val := r.URL.Query().Get("from_account_id"); val != "" {
		fromAccountId, err = strconv.Atoi(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if val := r.URL.Query().Get("to_account_id"); val != "" {
		toAccountId, err = strconv.Atoi(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if val := r.URL.Query().Get("pending"); val != "" {
		val, err := strconv.ParseBool(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
		pending = &val
	}

	if val := r.URL.Query().Get("accepted"); val != "" {
		val, err := strconv.ParseBool(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
		accepted = &val
	}
	fmt.Println(r.URL.Query().Get("pending"), r.URL.Query().Get("accepted"))

	filter := friend.NewQueryFilter()
	if fromAccountId != 0 {
		filter.WithFromAccountID(int64(fromAccountId))
	}
	if toAccountId != 0 {
		filter.WithToAccountID(int64(toAccountId))
	}
	if pending != nil {
		filter.WithPending(pending)
	}
	if accepted != nil {
		filter.WithAccepted(accepted)
	}

	if err = filter.Valid(); err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	limiter, err := limit.Parse(r, 1, 50)
	if err != nil {
		return err
	}

	data, err := h.core.ListRequests(ctx, filter, limiter)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, data, http.StatusOK)
}
