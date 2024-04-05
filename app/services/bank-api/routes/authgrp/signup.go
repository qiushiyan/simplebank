package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/user"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SignupRequest struct {
	Username string `json:"username" validate:"required,alphanum,username"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type SignupResponse struct {
	User userResponse `json:"user"`
}

// Signup godoc
//
//	@Summary		Signup
//	@Description	Signup with username, email and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SignupRequest	true	"request body"
//	@Success		201		{object}	SignupResponse
//	@Failure		400
//	@Failure		409
//	@Router			/signup [post]
func (h *Handler) Signup(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SignupRequest
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	user, err := h.core.Create(ctx, user.NewUser{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return err
	}

	response := SignupResponse{
		User: NewUserResponse(user),
	}

	return web.RespondJson(ctx, w, response, http.StatusCreated)
}
