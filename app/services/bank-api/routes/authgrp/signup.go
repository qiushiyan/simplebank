package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/user"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SignupRequest struct {
	Username string `json:"username" validate:"required,alphanum,username"`
	Email    string `json:"email"    validate:""`
	Password string `json:"password" validate:"required,password"`
}

type SignupResponse struct {
	User        userResponse `json:"user"`
	AccessToken string       `json:"access_token"`
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

	u, err := h.user.Create(ctx, user.NewUser{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return err
	}

	token, err := h.user.CreateToken(ctx, u, user.NewToken{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	if req.Email != "" {
		emailPayload := taskcommon.NewEmailDeliveryPayload(
			u.Username,
			"Welcome to SimpleBank",
			"signup-welcome",
		)

		h.task.CreateTask(taskcommon.TypeEmailDelivery, emailPayload)
	}

	response := SignupResponse{
		User:        NewUserResponse(u),
		AccessToken: token.Value,
	}

	return web.RespondJson(ctx, w, response, http.StatusCreated)
}
