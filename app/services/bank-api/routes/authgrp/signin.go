package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/user"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SigninRequest struct {
	//
	Username string `json:"username" validate:"required,alphanum,username"`
	Password string `json:"password" validate:"required,password"`
}

type SigninResponse struct {
	User        userResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

// Signin godoc
//
//	@Summary		Signin
//	@Description	Signin with username and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SigninRequest	true	"request body"
//	@Success		200		{object}	SigninResponse
//	@Failure		400
//	@Failure		409
//	@Router			/signin [post]
func (h *Handler) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SigninRequest
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	u, err := h.user.QueryByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	token, err := h.user.CreateToken(ctx, u, user.NewToken(req))
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, SigninResponse{
		User:        NewUserResponse(u),
		AccessToken: token.Value,
	}, http.StatusOK)
}
