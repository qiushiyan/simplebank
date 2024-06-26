package friendgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/friend"
	"github.com/qiushiyan/simplebank/business/data/limit"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type ListFriendshipResponse struct {
	Data []db_generated.Friendship `json:"data"`
}

// ListFriendship godoc
//
//	@Summary		List friendships
//	@Description	List friendship requests
//	@Tags			Friendship
//	@Accept			json
//	@Produce		json
//	@Param			from_account_id	query	int		false	"from account id"
//	@Param			to_account_id	query	int		false	"to account id"
//	@Param			pending			query	bool	false	"pending"
//	@Param			accepted		query	bool	false	"accepted"
//	@Param			page_id			query	int		false	"page id, default to 1"
//	@Param			page_size		query	int		false	"page size, default to 50"
//	@Security		Bearer
//	@Success		200	{object}	ListFriendshipResponse
//	@Failure		400
//	@Failure		409
//	@Router			/friend/list [get]
func (h *Handler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var fromAccountId int64
	var toAccountId int64
	var status *friend.Status
	var err error

	values := r.URL.Query()

	if val := values.Get("from_account_id"); val != "" {
		fromAccountId, err = cast.ToInt64E(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if val := values.Get("to_account_id"); val != "" {
		toAccountId, err = cast.ToInt64E(val)
		if err != nil {
			return web.NewError(err, http.StatusBadRequest)
		}
	}

	if val := values.Get("status"); val != "" {
		s, err := friend.ParseStatus(val)
		if err != nil {
			return web.NewError(
				friend.InvalidStatusError,
				http.StatusBadRequest,
			)
		}
		status = &s
	}

	filter := friend.NewQueryFilter()
	if fromAccountId != 0 {
		filter.WithFromAccountID(int64(fromAccountId))
	}
	if toAccountId != 0 {
		filter.WithToAccountID(int64(toAccountId))
	}
	if status != nil {
		filter.WithStatus(*status)
	}

	if err = filter.Valid(); err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	username := auth.GetUsername(ctx)
	// check from/to account id is owned by the user
	if fromAccountId != 0 {
		account, err := h.account.QueryById(ctx, int64(fromAccountId))
		if err != nil {
			return err
		}
		if account.Owner != username {
			return auth.NewForbiddenError(username)
		}
	}

	if toAccountId != 0 {
		account, err := h.account.QueryById(ctx, int64(toAccountId))
		if err != nil {
			return err
		}
		if account.Owner != username {
			return auth.NewForbiddenError(username)
		}
	}

	limiter, err := limit.Parse(r, 1, 50)
	if err != nil {
		return err
	}

	data, err := h.friend.ListRequests(ctx, filter, limiter)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, data, http.StatusOK)
}
