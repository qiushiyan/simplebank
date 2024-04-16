package friendgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/friend"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type UpdateRequest struct {
	// one of "pending", "accepted" or "rejected"
	Status string `json:"status" validate:"required"`
}

type UpdateFriendshipResponse struct {
	Data db_generated.Friendship `json:"data"`
}

// UpdateFriendship godoc
//
//	@Summary		Update friendship status
//	@Description	Update friendship status to be one of "pending", "accepted" or "rejected"
//	@Tags			Friendship
//
//	@Accept			json
//	@Produce		json
//	@Param			body	body	UpdateRequest	true	"request body"
//	@Param			id		path	int				true	"friendship id"
//	@Security		Bearer
//
//	@Success		200	{object}	UpdateFriendshipResponse
//	@Failure		400
//	@Failure		409
//	@Router			/friend/{id} [patch]
func (h *Handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req UpdateRequest
	if err := web.Decode(r, &req); err != nil {
		return err
	}

	fid, err := cast.ToInt64E(web.Param(r, "id"))
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	status, err := friend.ParseStatus(req.Status)
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	username := auth.GetUsername(ctx)
	// check that the to account id is owned by the user
	record, err := h.friend.GetFriendRequest(ctx, fid)
	if err != nil {
		return err
	}
	account, err := h.account.QueryById(ctx, record.ToAccountID)
	if err != nil {
		return err
	}

	if account.Owner != username {
		return auth.NewForbiddenError(username)
	}

	friend, err := h.friend.UpdateFriendRequest(ctx, fid, status)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, friend, http.StatusOK)
}
