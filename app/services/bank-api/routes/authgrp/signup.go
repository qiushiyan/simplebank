package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/core/user"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SignupRequest struct {
	// 3 to 20 characters, only letters and numbers
	Username string `json:"username" validate:"required,alphanum,username"`
	// Email address (Optional)
	Email string `json:"email"    validate:""`
	// 6 to 20 characters
	Password string `json:"password" validate:"required,password"`
}

type SignupResponse struct {
	// The created user model
	User userResponse `json:"user"`
	// Access token for the user
	AccessToken string `json:"access_token"`
	// Task id for the email delivery task, if email is provided in the request
	TaskId string `json:"task_id,omitempty"`
}

// Signup godoc
//
//	@Summary		Signup
//	@Description	Signup with username, email (optional) and password.
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
	var res SignupResponse
	err := web.Decode(r, &req)
	if err != nil {
		return err
	}

	result, err := h.user.CreateTx(ctx, user.NewUser{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}, func(user db_generated.User) (any, error) {
		return "", nil
	})

	if err != nil {
		return err
	}

	token, err := h.user.CreateToken(ctx, result.User, user.NewToken{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	res.User = NewUserResponse(result.User)
	res.AccessToken = token.Value
	res.TaskId = result.AfterCreateResult.(string)
	return web.RespondJson(ctx, w, res, http.StatusCreated)
}
