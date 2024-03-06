package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/user"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

const (
	adminUsername = "admin"
)

type SigninRequest struct {
	Username string `json:"username" validate:"required,alphanum,username"`
	Password string `json:"password" validate:"required,password"`
}

type SigninResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

func (h *Handler) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SigninRequest
	err := web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	if err = validate.Check(req); err != nil {
		return err
	}

	u, err := h.core.QueryByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	token, err := h.core.CreateToken(ctx, u, user.NewToken{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, SigninResponse{
		User:        NewUserResponse(u),
		AccessToken: token.GetToken(),
	}, http.StatusOK)
}
