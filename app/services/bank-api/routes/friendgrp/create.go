package friendgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/friend"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type CreateFriendRequest struct {
	// The sender account ID
	FromAccountId int64
	// The receiver account ID
	ToAccountId int64
}

type CreateFriendResponse struct {
	Data db_generated.Friendship `json:"data"`
}

// CreateFriendship godoc
//
//	@Summary		Create a friendship
//	@Description	Create a friendship between two accounts
//	@Tags			Friendship
//	@Accept			json
//	@Produce		json
//	@Param			body	body	CreateFriendRequest	true	"request body"
//	@Security		Bearer
//
//	@Success		201	{object}	CreateFriendResponse
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/friend/create [post]
func (h *Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req CreateFriendRequest
	if err := web.Decode(r, &req); err != nil {
		return err
	}

	// check from account is owned by the user
	username := auth.GetUsername(ctx)
	account, err := h.account.QueryById(ctx, req.FromAccountId)
	if err != nil {
		return err
	}
	if account.Owner != username {
		return auth.NewForbiddenError(username)
	}

	friendship, err := h.friend.NewRequest(ctx, friend.NewFriend(req))
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, friendship, http.StatusCreated)
}
