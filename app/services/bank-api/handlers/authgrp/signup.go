package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SignupRequest struct {
	Username string `json:"username" validate:"required,alphanum,username"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type SignupResponse struct {
	User UserResponse `json:"user"`
}

func (h *Handler) Signup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SignupRequest
	err := web.ParseBody(r, &req)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	if err := validate.Check(req); err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}

	user, err := h.store.CreateUser(ctx, db_generated.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	})

	if err != nil {
		return db.NewError(err)
	}

	response := SignupResponse{
		User: NewUserResponse(user),
	}

	return web.RespondJson(ctx, w, response, http.StatusCreated)
}
