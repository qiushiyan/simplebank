package friendgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/friend"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type CreateFriendRequest struct {
	FromAccountId int64
	ToAccountId   int64
}

func (h *Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateFriendRequest
	if err := web.Decode(r, &req); err != nil {
		return err
	}

	friendship, err := h.core.NewRequest(ctx, friend.NewFriend(req))
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, friendship, http.StatusCreated)
}
