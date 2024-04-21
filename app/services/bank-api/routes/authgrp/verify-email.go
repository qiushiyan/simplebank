package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/user"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type VerifyEmailRequest struct {
	Id   string `json:"id"   validate:"required"`
	Code string `json:"code" validate:"required"`
}

type VerifyEmailResponse struct {
	Ok bool `json:"ok"`
}

// VerifyEmail godoc
//
//	@Summary		Verify email
//	@Description	Verify email with the id and code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			body	body		VerifyEmailRequest	true	"request body"
//	@Success		200		{object}	VerifyEmailResponse
//	@Failure		400
//	@Failure		404
//	@Router			/verify-email [post]
func (h *Handler) VerifyEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req VerifyEmailRequest
	if err := web.Decode(r, &req); err != nil {
		return err
	}

	id, err := cast.ToInt64E(req.Id)
	if err != nil {
		return web.NewError(err, http.StatusBadRequest)
	}

	username := auth.GetUsername(ctx)
	u, err := h.user.QueryByUsername(ctx, username)
	if err != nil {
		return err
	}

	_, err = h.user.VerifyEmail(ctx, u, user.FinishVerifyEmail{
		Id:   id,
		Code: req.Code,
	})

	if err != nil {
		return err
	}

	return web.RespondJsonPlain(ctx, w, VerifyEmailResponse{Ok: true}, http.StatusOK)
}
