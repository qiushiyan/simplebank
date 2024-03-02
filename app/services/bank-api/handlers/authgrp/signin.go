package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/auth/token"
	db "github.com/qiushiyan/simplebank/business/db/core"
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

	user, err := h.store.GetUser(ctx, req.Username)
	if err != nil {
		return db.NewError(err)
	}

	if !auth.VerifyPassword(user.HashedPassword, req.Password) {
		return auth.NewAuthError("incorrect password")
	}

	roles := []string{"USER"}
	if user.Username == adminUsername {
		roles = append(roles, "ADMIN")
	}
	t, err := token.NewToken(user.Username, roles, 0)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}

	return web.RespondJson(ctx, w, SigninResponse{
		User:        NewUserResponse(user),
		AccessToken: t.GetToken(),
	}, http.StatusOK)
}
